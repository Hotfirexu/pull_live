package base

import (
	"bufio"
	"fmt"
	"github.com/q191201771/naza/pkg/bininfo"
	"github.com/q191201771/naza/pkg/nazalog"
	"os"
	"runtime"
	"strings"
	"time"
)

var Log = nazalog.GetGlobalLogger()

var startTime string

func GetWd() string {
	dir, _ := os.Getwd()
	return dir
}

func ReadableNowTime() string {
	return time.Now().Format("2006-01-02 15:04:05.999")
}

func LogoutStartInfo() {
	Log.Infof("     start: %s", startTime)
	Log.Infof("        wd: %s", GetWd())
	Log.Infof("      args: %s", strings.Join(os.Args, " "))
	Log.Infof("   bininfo: %s", bininfo.StringifySingleLine())
	Log.Infof("   version: %s", LalFullInfo)
	Log.Infof("    github: %s", LalGithubSite)
	Log.Infof("       doc: %s", LalDocSite)
}

func init() {
	startTime = ReadableNowTime()
}

func OsExitAndWaitPressIfWindows(code int) {
	if runtime.GOOS == "windows" {
		_, _ = fmt.Fprintf(os.Stderr, "Press Enter to exit...")
		r := bufio.NewReader(os.Stdin)
		_, _ = r.ReadByte()
	}
	os.Exit(code)
}
