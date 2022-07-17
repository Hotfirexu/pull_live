package httpflv

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/q191201771/naza/pkg/nazalog"
	"net"
	"net/http"
	"pull_live/base"
	thttp "pull_live/http"
	"pull_live/pkg/connection"
	"sync"
	"time"
)

type PullSessionOption struct {
	PullTimeoutMs int
	ReadTimeoutMs int
}

var defaultPullSessionOption = PullSessionOption{
	PullTimeoutMs: 10000,
	ReadTimeoutMs: 0,
}

type PullSession struct {
	option      PullSessionOption
	conn        connection.Connection
	sessionStat base.BasicSessionStat
	urlCtx      base.UrlContext
	disposeOnce sync.Once
}

type ModPullSessionOption func(option *PullSessionOption)

func NewPullSession(modOptions ...ModPullSessionOption) *PullSession {
	option := defaultPullSessionOption
	for _, fn := range modOptions {
		fn(&option)
	}

	s := &PullSession{
		option:      option,
		sessionStat: base.NewBasicSessionStat(base.SessionTypeFlvPull, ""),
	}

	return s
}

type OnReadFlvTag func(tag Tag)

func (session *PullSession) Pull(rawUrl string, onReadFlvTag OnReadFlvTag) error {
	fmt.Println("key=", session.UniqueKey(), " url=", rawUrl)

	var (
		ctx    context.Context
		cancel context.CancelFunc
	)

	if session.option.PullTimeoutMs == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), time.Duration(session.option.PullTimeoutMs)*time.Millisecond)
	}

	defer cancel()
	return session.pullContext(ctx, rawUrl, onReadFlvTag)
}

func (session *PullSession) pullContext(ctx context.Context, rawUrl string, onReadFlvTag OnReadFlvTag) error {
	errChan := make(chan error, 1)
	url := rawUrl

	// 异步握手
	go func() {
		for {
			if err := session.connect(url); err != nil {
				errChan <- err
				return
			}

			if err := session.writeHttpRequest(); err != nil {
				errChan <- err
				return
			}

			statusCode, headers, err := session.readHttpResponse()
			if err != nil {
				errChan <- err
				return
			}

			// 处理跳转
			if statusCode == "302" || statusCode == "301" {
				url = headers.Get("Location")
				if url == "" {
					nazalog.Warnf("[%s] redirect but Location not found. headers = %+v",
						session.UniqueKey(), headers)
					errChan <- nil
					return
				}
				_ = session.conn.Close()
				nazalog.Debugf("[%s] redirect to %s", session.UniqueKey(), url)
				continue
			}

			errChan <- nil
			return
		}
	}()

	// 等待握手结束，或者超时通知
	select {
	case <-ctx.Done():
		// 注意： 如果超时了，可能连接已经建立，这里一定要dispose，避免泄露
		_ = session.dispose(nil)
		return ctx.Err()
	case err := <-errChan:
		// 握手消息，不为nil则握手失败
		if err != nil {
			_ = session.dispose(err)
			return err
		}
	}

	// 握手成功，开始收数据协程
	go session.runReadLoop(onReadFlvTag)
	return nil

}

func (session *PullSession) UniqueKey() string {
	return session.sessionStat.UniqueKey()
}

func (session *PullSession) connect(rawUrl string) (err error) {
	session.urlCtx, err = base.ParseHttpFlvUrl(rawUrl)
	if err != nil {
		return err
	}

	session.sessionStat.SetRemoteAddr(session.urlCtx.HostWithPort)
	fmt.Println("key=", session.UniqueKey(), " > tcp connect. ", session.urlCtx.HostWithPort)

	var conn net.Conn
	if session.urlCtx.Scheme == "https" {
		conf := &tls.Config{
			InsecureSkipVerify: true,
		}

		conn, err = tls.Dial("tcp", session.urlCtx.HostWithPort, conf)
	} else {
		conn, err = net.Dial("tcp", session.urlCtx.HostWithPort)
	}
	if err != nil {
		return err
	}
	fmt.Println("key=", session.UniqueKey(), "tcp connect succ. remote=", conn.RemoteAddr().String())

	session.conn = connection.New(conn, func(option *connection.Option) {
		option.ReadBufSize = readBufSize
		option.WriteTimeoutMs = session.option.ReadTimeoutMs
		option.ReadTimeoutMs = session.option.ReadTimeoutMs
	})

	return nil
}

func (session *PullSession) writeHttpRequest() error {
	// 发送http GET请求
	req := fmt.Sprintf("GET %s HTTP/1.0\r\nUser-Agent: %s\r\nAccept: */*\r\n"+
		"Range: byte=0-\r\nConnection: close\r\nHost: %s\r\nIcy-MetaData: 1\r\n",
		session.urlCtx.PathWithRawQuery, base.HttpflvSubSessionServer, session.urlCtx.StdHost)
	_, err := session.conn.Write([]byte(req))
	return err
}

func (session *PullSession) readHttpResponse() (statusCode string, headers http.Header, err error) {
	var statusLine string
	if statusLine, headers, err = thttp.ReadHttpHeader(session.conn); err != nil {
		return
	}
	if _, statusCode, _, err = thttp.ParseHttpStatusLine(statusLine); err != nil {
		return
	}
	nazalog.Debugf("[%s] < R http response header. statusLine= %s", session.UniqueKey(), statusLine)
	return
}

func (session *PullSession) dispose(err error) error {
	var retErr error
	session.disposeOnce.Do(func() {
		nazalog.Infof("[%s lifecycle dispose http-flv pull session. err: %+v]",
			session.UniqueKey(), err)
		if session.conn == nil {
			retErr = base.ErrSessionNotStarted
			return
		}
		retErr = session.conn.Close()
	})

	return retErr
}

func (session *PullSession) runReadLoop(onReadFlvTag OnReadFlvTag) {
	var err error
	defer func() {
		_ = session.dispose(err)
	}()

	if _, err = session.readFlvHeader(); err != nil {
		return
	}

	for {
		var tag Tag
		if tag, err = session.readTag(); err != nil {
			return
		}

		onReadFlvTag(tag)
	}
}

func (session *PullSession) readFlvHeader() ([]byte, error) {
	flvHeader := make([]byte, flvHeaderSize)
	if _, err := session.conn.ReadAtLeast(flvHeader, flvHeaderSize); err != nil {
		return flvHeader, err
	}

	nazalog.Debugf("[%s] < R http flv header", session.UniqueKey())

	// todo: check value is valid
	return flvHeader, nil
}

func (session *PullSession) readTag() (Tag, error) {
	return ReadTag(session.conn)
}

func (session *PullSession) GetStat() base.StatSession {
	return session.sessionStat.GetStatWithConn(session.conn)
}

func (session *PullSession) WaitChan() <-chan error {
	return session.conn.Done()
}
