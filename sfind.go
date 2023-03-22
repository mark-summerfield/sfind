// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	_ "embed"
	"sync"
)

//go:embed Version.dat
var Version string

func main() {
	config := getConfig()
	var wg sync.WaitGroup
	for i := 0; i < len(config.paths); i++ {
		wg.Add(1)
		i := i
		go func() {
			defer wg.Done()
			processPath(i, config)
		}()
	}
	wg.Wait()
}
