package main

import (
	"os"
	"runtime/trace"
	"sync"
)

func main() {
	var (
		wg sync.WaitGroup
		me sync.Mutex
	)

	trace.Start(os.Stderr)
	defer trace.Stop()

	mp := make(map[int]int)

	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			for j := 0; j < 1000000; j++ {
				me.Lock()
				mp[num] = num
				me.Unlock()
			}
		}(i)
	}

	wg.Wait()

}
