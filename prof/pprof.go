package prof

import (
	"os"
	"runtime/pprof"
	"time"
)

var hc int
var cc int

func Heap() error {
	heap, err := os.Create("sharp-heap-profile")
	if err != nil {
		return err
	}
	defer heap.Close()

	hc++
	return pprof.WriteHeapProfile(heap)
}

func CPU() error {

	cpu, err := os.Create("sharp-cpu-profile")
	if err != nil {
		return err
	}

	defer cpu.Close()
	cc++

	err = pprof.StartCPUProfile(cpu)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 30)
	pprof.StopCPUProfile()

	return nil
}

func Monitor(t time.Duration) {
	go func() {
		tick := time.Tick(t)
		for {
			select {
			case <-tick:
				CPU()
				Heap()
			}
		}
	}()
}
