//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Create a process
	proc := MockProcess{}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		// Run the process (blocking)
		proc.Run()
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	go func() {
		sigCount := 0
		for {
			<-signalCh
			sigCount++
			if sigCount == 1 {
				go proc.Stop()
			} else if sigCount == 2 {
				os.Exit(1)
			}
		}
	}()

	wg.Wait()
}
