package version

import (
	"fmt"
	"runtime"
)

var (
	Version   string = ""
	BuildTime string = ""
	SoftName  string = ""
)

func PrintVersion() string {
	return fmt.Sprintf("%s %s %s/%s, %s, %s\n",
		SoftName,
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.Version(),
		BuildTime,
	)
}
