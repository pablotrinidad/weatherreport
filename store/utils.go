package store

import "sync"

// callConcurrent perform fns concurrently.
func callConcurrent(fns []func()) {
	var wg sync.WaitGroup
	for i := range fns {
		f := fns[i]
		wg.Add(1)
		go func() {
			defer wg.Done()
			f()
		}()
	}
	wg.Wait()
}