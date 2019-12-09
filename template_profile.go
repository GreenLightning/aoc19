import (
	"os/signal"
	"runtime/pprof"
	"syscall"
)

var stopProfiling = make(chan bool)

// Usage: go profile()
func profile() {
	file, err := os.Create("today.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(file)

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-channel:
	case <-stopProfiling:
	}

	pprof.StopCPUProfile()
	os.Exit(0)
}
