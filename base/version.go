package base

import "strings"

const Version = "v0.0.1"

var (
	LalLibraryName = "lal"
	LalGithubRepo  = "github.com/HotFire"
	LalGithubSite  = "https://github.com/q191201771/lal"
	LalDocSite     = "https://pengrl.com/lal"
	LalFullInfo    = LalLibraryName + " " + Version + " (" + LalGithubRepo + ")"

	// LalVersionDot e.g. 0.12.3
	LalVersionDot string
)
var (
	// LalHttpflvSubSessionServer e.g. lal0.12.3
	HttpflvSubSessionServer string
)

func init() {
	LalVersionDot = strings.TrimPrefix(Version, "v")
	HttpflvSubSessionServer = LalLibraryName + LalVersionDot
}
