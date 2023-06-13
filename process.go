// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/mark-summerfield/gset"
)

func processPath(index int, config *config) {
	_ = filepath.WalkDir(config.paths[index],
		func(path string, de fs.DirEntry, err error) error {
			if err == nil {
				if de.IsDir() {
					if skipFolder(path, config) {
						return fs.SkipDir
					}
				} else if info, err := de.Info(); err == nil {
					if validFilename(path, info, config) {
						fmt.Println(path)
						return nil
					}
				}
			}
			return err
		})
}

func skipFolder(path string, config *config) bool {
	if config.casefold {
		path = strings.ToLower(path)
	}
	base := filepath.Base(path)
	if base != "." && strings.HasPrefix(base, ".") { // skip hidden
		return true
	}
	path = filepath.ToSlash(path)
	parts := gset.New(strings.Split(path, "/")...)
	for _, name := range config.excludes {
		if config.casefold {
			name = strings.ToLower(name)
		}
		if parts.Contains(name) {
			return true
		}
	}
	return false
}

func validFilename(path string, info fs.FileInfo, config *config) bool {
	if config.from.After(info.ModTime()) {
		return false
	}
	base := filepath.Base(path)
	if strings.HasPrefix(base, ".") { // skip hidden
		return false
	}
	if config.casefold {
		base = strings.ToLower(base)
	}
	for _, glob := range config.globs {
		if config.casefold {
			glob = strings.ToLower(glob)
		}
		if matched, _ := filepath.Match(glob, base); matched {
			return true
		}
	}
	return false
}
