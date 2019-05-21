/*
	Copyright 2017-2018 OneLedger

	A fullnode for the OneLedger chain. Includes cli arguments to initialize, restart, etc.
*/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

func main() {

	rs := rand.NewSource(time.Now().UnixNano())
	f, err := os.Create(fmt.Sprintf("/tmp/profile%d_cpu.pprof", rs.Int63()))
	if err != nil {
		fmt.Println(err, ": error opening file")
		return
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		fmt.Println(err, ": err starting cpu profiling")
		return
	}

	Execute() // Pass control to Cobra

	pprof.StopCPUProfile()
}
