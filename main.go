package main

import (
	"fmt"
	"github.com/doout/cps842/cmd"
	"runtime"
	"time"
)

func main() {
	cmd.Execute()

}

func timeIt(fn func()) {
	start := time.Now()
	fn()
	end := time.Now()
	fmt.Println(end.Sub(start))
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
