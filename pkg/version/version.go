package version

import (
	"fmt"
	"runtime"
)

var (
	Version   string
	BuildDate string
)

func PrintVersion() string {
	return fmt.Sprintf("version %s, built on %s, %s",
		Version,
		BuildDate,
		runtime.Version())
}
