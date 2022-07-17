package base

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultRtmpPort  = 1935
	DefaultHttpPort  = 80
	DefaultHttpsPort = 443
	DefaultRtspPort  = 554
	DefaultRtmpsPort = 443
)

type UrlPathContext struct {
	PathWithRawQuery    string
	Path                string
	PathWithoutLastItem string
	LastItemOfPath      string
	RawQuery            string
}

type UrlContext struct {
	Url string

	Scheme       string
	UserName     string
	Password     string
	StdHost      string
	HostWithPort string
	Host         string
	Port         int

	PathWithRawQuery    string
	Path                string
	PathWithoutLastItem string // 注意，没有前面的'/'，也没有后面的'/'
	LastItemOfPath      string // 注意，没有前面的'/'
	RawQuery            string

	RawUrlWithoutUserInfo string
	filenameWithoutType   string
	fileType              string
}

func ParseHttpFlvUrl(rawUrl string) (ctx UrlContext, err error) {
	return parseHttpUrl(rawUrl, ".flv")
}

func parseHttpUrl(rawUrl string, fileType string) (ctx UrlContext, err error) {
	ctx, err = ParseUrl(rawUrl, -1)
	if err != nil {
		return ctx, err
	}
	if (ctx.Scheme != "http" && ctx.Scheme != "https") || ctx.Host == "" || ctx.Path == "" ||
		!strings.HasSuffix(ctx.LastItemOfPath, fileType) {
		return ctx, fmt.Errorf("%w. url=%s", ErrInvalidUrl, rawUrl)
	}
	return
}

// ParseUrl
//
// @param defaultPort: 注意，如果rawUrl中显示指定了端口，则该参数不生效
//                     注意，如果设置为-1，内部依然会对常见协议(http, https, rtmp, rtsp)设置官方默认端口
//
func ParseUrl(rawUrl string, defaultPort int) (ctx UrlContext, err error) {
	ctx.Url = rawUrl

	stdUrl, err := url.Parse(rawUrl)
	if err != nil {
		return ctx, err
	}

	if stdUrl.Scheme == "" {
		return ctx, fmt.Errorf("%w. url=%s", ErrInvalidUrl, rawUrl)
	}

	if defaultPort == -1 {
		switch stdUrl.Scheme {
		case "http":
			defaultPort = DefaultHttpPort
		case "https":
			defaultPort = DefaultHttpsPort
		case "rtmp":
			defaultPort = DefaultRtmpPort
		case "rtsp":
			defaultPort = DefaultRtspPort

		}
	}

	ctx.Scheme = stdUrl.Scheme
	ctx.StdHost = stdUrl.Host
	ctx.UserName = stdUrl.User.Username()
	ctx.Password, _ = stdUrl.User.Password()

	h, p, err := net.SplitHostPort(stdUrl.Host)
	if err != nil {
		ctx.Host = stdUrl.Host
		if defaultPort == -1 {
			ctx.HostWithPort = stdUrl.Host
		} else {
			ctx.HostWithPort = net.JoinHostPort(stdUrl.Host, fmt.Sprintf("%d", defaultPort))
			ctx.Port = defaultPort
		}
	} else {
		ctx.Port, err = strconv.Atoi(p)
		if err != nil {
			return ctx, err
		}

		ctx.Host = h
		ctx.HostWithPort = stdUrl.Host
	}
	pathCtx, err := parseUrlPath(stdUrl)
	if err != nil {
		return ctx, err
	}

	ctx.PathWithRawQuery = pathCtx.PathWithRawQuery
	ctx.Path = pathCtx.Path
	ctx.PathWithoutLastItem = pathCtx.PathWithoutLastItem
	ctx.LastItemOfPath = pathCtx.LastItemOfPath
	ctx.RawQuery = pathCtx.RawQuery

	ctx.RawUrlWithoutUserInfo = fmt.Sprintf("%s//%s%s", ctx.Scheme, ctx.StdHost, ctx.PathWithRawQuery)
	return ctx, nil
}

func parseUrlPath(stdUrl *url.URL) (ctx UrlPathContext, err error) {
	ctx.Path = stdUrl.Path

	index := strings.LastIndexByte(ctx.Path, '/')
	if index == -1 {
		ctx.PathWithoutLastItem = ""
		ctx.LastItemOfPath = ""
	} else if index == 0 {
		if ctx.Path == "/" {
			ctx.PathWithoutLastItem = ""
			ctx.LastItemOfPath = ""
		} else {
			ctx.PathWithoutLastItem = ""
			ctx.LastItemOfPath = ctx.Path[1:]
		}
	} else {
		ctx.PathWithoutLastItem = ctx.Path[1:index]
		ctx.LastItemOfPath = ctx.Path[index+1:]
	}

	ctx.RawQuery = stdUrl.RawQuery

	if ctx.RawQuery == "" {
		ctx.PathWithRawQuery = ctx.Path
	} else {
		ctx.PathWithRawQuery = fmt.Sprintf("%s?%s", ctx.Path, ctx.RawQuery)
	}

	return ctx, nil
}