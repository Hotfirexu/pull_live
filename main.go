package main

import (
	obytes "bytes"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/q191201771/naza/pkg/nazalog"
	"pull_live/avc"
	"pull_live/base"
	"pull_live/httpflv"
	"pull_live/pkg/bele"
	"pull_live/pkg/birate"
	"pull_live/pkg/bytes"
	"pull_live/rtmp"
	"strconv"
	"strings"
	"time"
)

var (
	printStatFlag        = true
	printEveryTagFlag    = true
	printMetaData        = true
	timestampCheckFlag   = true
	analysisVideoTagFlag = true
)

var (
	prevAudioTs = int64(-1)
	prevVideoTs = int64(-1)
	prevTs      = int64(-1)
	prevIdrTs   = int64(-1)
	diffIdrTs   = int64(-1)
)

func main() {
	_ = nazalog.Init(func(option *nazalog.Option) {
		option.AssertBehavior = nazalog.AssertFatal
	})

	defer nazalog.Sync()
	base.LogoutStartInfo()

	url := parseFlag()

	session := httpflv.NewPullSession()

	brTotal := birate.New(func(option *birate.Option) {
		option.WindowMs = 5000
	})

	brAudio := birate.New(func(option *birate.Option) {
		option.WindowMs = 5000
	})

	brVideo := birate.New(func(option *birate.Option) {
		option.WindowMs = 5000
	})

	videoCtsNotZeroCount := 0

	go func() {
		for {
			time.Sleep(5 * time.Second)
			if printStatFlag {
				nazalog.Debugf("stat. total:%dKb/s, audio=%dKb/s, video=%dKb/s, videoCtsNotZeroCount=%d, diffIdrTs=%d",
					int(brTotal.Rate()), int(brAudio.Rate()), int(brVideo.Rate()), videoCtsNotZeroCount, diffIdrTs)
			}

		}
	}()

	err := session.Pull(url, func(tag httpflv.Tag) {
		if printEveryTagFlag {
			nazalog.Debugf("header=%+v, hex=%s", tag.Header, hex.Dump(bytes.Prefix(tag.Payload(), 32)))
		}

		brTotal.Add(len(tag.Raw))

		switch tag.Header.Type {
		case httpflv.TagTypeMetadata:
			if printMetaData {
				nazalog.Debugf("-------\n%s", hex.Dump(tag.Payload()))

				opa, err := rtmp.ParseMetaData(tag.Payload())
				nazalog.Assert(nil, err)
				var buf obytes.Buffer
				buf.WriteString(fmt.Sprintf("------\ncount:%d\n", len(opa)))
				for _, op := range opa {
					buf.WriteString(fmt.Sprintf(" %s: %+v\n", op.Key, op.Value))
				}
				nazalog.Debugf("%+v", buf.String())
			}

		case httpflv.TagTypeAudio:
			nazalog.Debugf("header=%+v, %+v", tag.Header, tag.IsAacSeqHeader())
			brAudio.Add(len(tag.Raw))

			if tag.IsAacSeqHeader() {
				s := session.GetStat()
				nazalog.Infof("aac seq header. readBytes=%d, %s", s.ReadBytesSum,
					hex.EncodeToString(tag.Payload()))
			}

			if timestampCheckFlag {
				if prevAudioTs != -1 && int64(tag.Header.TimeStamp) < prevAudioTs {
					nazalog.Errorf("audio timestamp error, less than prev audio timestamp. "+
						"header=%+v, prevAudioTs=%d, diff=%d", tag.Header, prevAudioTs,
						int64(tag.Header.TimeStamp)-prevAudioTs)
				}
				if prevTs != -1 && int64(tag.Header.TimeStamp) < prevTs {
					nazalog.Warnf("audio timestamp error. less than prev global timestamp. header=%+v, prevTs=%d, diff=%d", tag.Header, prevTs, int64(tag.Header.TimeStamp)-prevTs)
				}
			}

			prevAudioTs = int64(tag.Header.TimeStamp)
			prevTs = int64(tag.Header.TimeStamp)
		case httpflv.TagTypeVideo:
			analysisVideoTag(tag)

			videoCts := bele.BeUint24(tag.Raw[13:])
			if videoCts != 0 {
				videoCtsNotZeroCount++
			}

			brVideo.Add(len(tag.Raw))

			if timestampCheckFlag {
				if prevVideoTs != -1 && int64(tag.Header.TimeStamp) < prevVideoTs {
					nazalog.Errorf("video timestamp error, less than prev video timestamp. header=%+v, prevVideoTs=%d, diff=%d", tag.Header, prevVideoTs, int64(tag.Header.TimeStamp)-prevVideoTs)
				}
				if prevTs != -1 && int64(tag.Header.TimeStamp) < prevTs {
					nazalog.Warnf("video timestamp error, less than prev global timestamp. header=%+v, prevTs=%d, diff=%d", tag.Header, prevTs, int64(tag.Header.TimeStamp)-prevTs)
				}
			}
			prevVideoTs = int64(tag.Header.TimeStamp)
			prevTs = int64(tag.Header.TimeStamp)

		}

	})

	nazalog.Assert(nil, err)
	err = <-session.WaitChan()

	nazalog.Errorf("< session.WaitChan err=%+v", err)

}

const (
	typeUnknown uint8 = 1
	typeAvc     uint8 = 2
	typeHevc    uint8 = 3
)

var t uint8 = typeUnknown

func analysisVideoTag(tag httpflv.Tag) {
	var buf obytes.Buffer
	if tag.IsVideoKeySeqHeader() {
		if tag.IsAvcKeySeqHeader() {
			t = typeAvc
			buf.WriteString(" [AVC SeqHeader] ")
			sps, pps, err := avc.ParseSpsPpsFromSeqHeader(tag.Payload())
			if err != nil {
				buf.WriteString(" parse sps pps failed.")
			}
			nazalog.Debugf("sps:%s, pps:%s", hex.Dump(sps), hex.Dump(pps))
		} //else if tag.IsHevcKeySeqHeader() {
		//t = typeHevc
		//buf.WriteString(" [HEVC SeqHeader] ")
		//buf.WriteString(hex.Dump(tag.Payload()))
		//if _, _, _, err := hevc.ParseVpsSpsPpsFromSeqHeader(tag.Payload()); err != nil {
		//	buf.WriteString(" parse vps sps pps failed.")
		//}
		//}
	} else {
		cts := bele.BeUint24(tag.Payload()[2:])
		buf.WriteString(fmt.Sprintf("%+v, cts=%d, pts=%d", tag.Header, cts, tag.Header.TimeStamp+cts))

		body := tag.Payload()[5:]
		nals, err := avc.SplitNaluAvcc(body)
		nazalog.Assert(nil, err)

		for _, nal := range nals {
			switch t {
			case typeAvc:
				if avc.ParseNaluType(nal[0]) == avc.NaluTypeIdrSlice {
					nazalog.Debugf("IDR:%s", hex.Dump(bytes.Prefix(nal, 128)))
					if prevIdrTs != int64(-1) {
						diffIdrTs = int64(tag.Header.TimeStamp) - prevIdrTs
					}
					prevIdrTs = int64(tag.Header.TimeStamp)
				}
				if avc.ParseNaluType(nal[0]) == avc.NaluTypeSei {
					delay := SeiDelayMs(nal)
					if delay != -1 {
						buf.WriteString(fmt.Sprintf("delay: %dms", delay))
					}
				}
				sliceTypeReadable, _ := avc.ParseSliceTypeReadable(nal)
				buf.WriteString(fmt.Sprintf(" [%s(%s)(%d)] ", avc.ParseNaluTypeReadable(nal[0]), sliceTypeReadable, len(nal)))
			case typeHevc:
				//if hevc.ParseNaluType(nal[0]) == hevc.NaluTypeSei {
				//	delay := SeiDelayMs(nal)
				//	if delay != -1 {
				//		buf.WriteString(fmt.Sprintf("delay: %dms", delay))
				//	}
				//}
				//buf.WriteString(fmt.Sprintf(" [%s(%d)] ", hevc.ParseNaluTypeReadable(nal[0]), nal[0]))
			}
		}
	}
	if analysisVideoTagFlag {
		nazalog.Debug(buf.String())
	}
}

// 注意，SEI的内容是自定义格式，解析的代码不具有通用性
func SeiDelayMs(seiNalu []byte) int {
	//nazalog.Debugf("sei: %s", hex.Dump(seiNalu))
	items := strings.Split(string(seiNalu), ":")
	if len(items) != 3 {
		return -1
	}

	a, err := strconv.ParseInt(items[1], 10, 64)
	if err != nil {
		return -1
	}
	t := time.Unix(a/1e3, a%1e3)
	d := time.Now().Sub(t)
	return int(d.Nanoseconds() / 1e6)
}

func parseFlag() string {
	url := flag.String("i", "", "http-flv url")
	flag.Parse()

	if *url == "" {
		flag.Usage()
		base.OsExitAndWaitPressIfWindows(1)
	}
	return *url
}
