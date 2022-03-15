package version

import (
	"fmt"
)

var (
	Version   = "dev"
	Commit    = "n/a"
	BuildDate = "n/a"
)

func BuildVersion(appName string) string {
	return fmt.Sprintf("%s version %s\nCommit: %s\nBuildDate: %s", appName, Version, Commit, BuildDate)
}
