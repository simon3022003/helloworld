package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

// RunCPULoad run CPU load in specify cores count and percentage
func RunCPULoad() {
	coresCount := (runtime.NumCPU())
	runtime.GOMAXPROCS(coresCount)
	done := make(chan int)
	for i := 0; i < coresCount; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
			}
		}()

	}
	time.Sleep(time.Second * 300)
	close(done)
}

// Load simulates a high load event for 5 min
func Load(w http.ResponseWriter, r *http.Request) {
	go RunCPULoad()
	fmt.Println("Generating load for 5 min")
	w.WriteHeader(http.StatusOK)
	return
}
