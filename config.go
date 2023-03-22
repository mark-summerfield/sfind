// Copyright © 2023 Mark Summerfield. All rights reserved.
// License: GPL-3

package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mark-summerfield/clip"
)

func getConfig() *config {
	config := config{from: time.UnixMilli(0)}
	parser := clip.NewParserVersion(Version)
	parser.PositionalHelp = "Paths to search [default: .]"
	parser.PositionalCount = clip.ZeroOrMorePositionals
	_ = parser.SetPositionalVarName("PATH")
	parser.LongDesc = "Searches for matching files"
	fromOpt := parser.Str("from", "The earliest date to search from. "+
		"Can use 'today' or 'yesterday' or an int (up to that many days "+
		"ago), or an ISO8601 format date (e.g., 2023-05-22) [default: "+
		"any date]", "")
	globsOpt := parser.Strs("glob", "The file globs to match (e.g., "+
		"'*.py' '*.pyw') [default: any file]")
	excludeOpt := parser.Strs("exclude", "Paths to exclude [default: none]")
	excludeOpt.SetShortName('x')
	if err := parser.Parse(); err != nil {
		parser.OnError(err)
	}
	if fromOpt.Given() {
		from, err := parseFrom(fromOpt.Value())
		if err != nil {
			parser.OnError(err)
		} else {
			config.from = from
		}
	}
	if globsOpt.Given() {
		config.globs = globsOpt.Value()
		for _, glob := range config.globs {
			_, err := filepath.Match(glob, "")
			if err != nil {
				parser.OnError(fmt.Errorf("glob pattern: %q: %w", glob,
					err))
			}
		}
	} else {
		config.globs = []string{"*"}
	}
	if excludeOpt.Given() {
		config.excludes = excludeOpt.Value()
	} else {
		config.excludes = make([]string, 0)
	}
	if len(parser.Positionals) > 0 {
		config.paths = parser.Positionals
	} else {
		config.paths = []string{"."}
	}
	return &config
}

func parseFrom(text string) (time.Time, error) {
	now := time.Now()
	zero := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0,
		now.Location())
	text = strings.ToLower(text)
	if text == "today" {
		return zero, nil
	}
	if text == "yesterday" {
		return zero.AddDate(0, 0, -1), nil
	}
	i, err := strconv.Atoi(text)
	if err == nil {
		return zero.AddDate(0, 0, -i), nil
	}
	return time.ParseInLocation(time.DateOnly, text, now.Location())
}

type config struct {
	from     time.Time
	globs    []string
	excludes []string
	paths    []string
}

func (me config) String() string {
	return fmt.Sprintf("from=%s\nglobs=[%s]\nexcludes=[%s]\npaths=[%s]",
		me.from, strings.Join(me.globs, " "),
		strings.Join(me.excludes, " "), strings.Join(me.paths, " "))
}
