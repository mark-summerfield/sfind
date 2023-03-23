// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	_ "embed"
	"fmt"
	"sync"
)

//go:embed Version.dat
var Version string

func main() {
	config := getConfig()
	if config.debug {
		fmt.Println(config)
	} else {
		var wg sync.WaitGroup
		for i := 0; i < len(config.paths); i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				processPath(index, config)
			}(i)
		}
		wg.Wait()
	}
}
