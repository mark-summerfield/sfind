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
	const comma = ","
	parser := makeParser()
	suffixesOpt := parser.Str("suffixes", suffixesDesc, "")
	containsOpt := parser.Str("contains", containsDesc, "")
	globsOpt := parser.Str("glob", globsDesc, "")
	fromOpt := parser.Str("from", fromDesc, "")
	excludeOpt := parser.Str("exclude", excludeDesc, "")
	excludeOpt.SetShortName('x')
	debugOpt := parser.Flag("debug", "print config and quit")
	debugOpt.Hide()
	if err := parser.Parse(); err != nil {
		parser.OnError(err)
	}
	config := &config{from: time.UnixMilli(0)}
	if fromOpt.Given() {
		updateFrom(parser, fromOpt.Value(), config)
	}
	globs := make([]string, 0)
	if globsOpt.Given() {
		globs = strings.Split(globsOpt.Value(), comma)
	}
	if containsOpt.Given() {
		for _, contains := range strings.Split(containsOpt.Value(), comma) {
			globs = append(globs, fmt.Sprintf("*%s*", contains))
		}
	}
	suffixes := make([]string, 0)
	if suffixesOpt.Given() {
		suffixes = strings.Split(suffixesOpt.Value(), comma)
	}
	updateGlobs(parser, globs, suffixes, config)
	if excludeOpt.Given() {
		config.excludes = strings.Split(excludeOpt.Value(), comma)
	}
	if len(parser.Positionals) > 0 {
		config.paths = parser.Positionals
	} else {
		config.paths = []string{"."}
	}
	config.debug = debugOpt.Value()
	return config
}

func makeParser() *clip.Parser {
	parser := clip.NewParserVersion(Version)
	parser.PositionalHelp = "Paths to search [default: .]"
	parser.PositionalCount = clip.ZeroOrMorePositionals
	parser.MustSetPositionalVarName("PATH")
	parser.LongDesc = longDesc
	parser.EndDesc = endDesc
	return &parser
}

func updateFrom(parser *clip.Parser, text string, config *config) {
	now := time.Now()
	base := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0,
		now.Location())
	text = strings.ToLower(text)
	if text == "today" {
		config.from = base
	} else if text == "yesterday" {
		config.from = base.AddDate(0, 0, -1)
	} else {
		i, err := strconv.Atoi(text)
		if err == nil {
			config.from = base.AddDate(0, 0, -i)
		} else {
			base, err := time.ParseInLocation(time.DateOnly, text,
				now.Location())
			if err == nil {
				config.from = base
			} else {
				parser.OnError(err)
			}
		}
	}
}

func updateGlobs(parser *clip.Parser, globs, suffixes []string,
	config *config) {
	if len(suffixes) > 0 {
		for _, suffix := range suffixes {
			globs = append(globs, fmt.Sprintf("*.%s", suffix))
		}
	}
	if len(globs) > 0 {
		config.globs = globs
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
}

type config struct {
	debug    bool
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

const (
	longDesc = "Searches for matching files."
	endDesc  = "Examples `sfind --from today` find files changed " +
		"today in the current folder and its subfolders; " +
		"`sfind -f0 -s go` find Go files changed today; " +
		"`-sfind -f1 -s py,pyw ~/app` find Python files changed since " +
		"yesterday in the $HOME/app folder."
	fromDesc = "The earliest date to search from. " +
		"Can use 'today' (or 0) or 'yesterday' (or 1) or an int (up " +
		"to that many days ago), or an ISO8601 format date " +
		"(e.g., 2023-05-22) [default: any date]."
	suffixesDesc = "The comma-separated file suffixes to match (e.g., " +
		"py,pyw) [default: any file]."
	globsDesc = "The comma-separated file globs to match (e.g., " +
		"'*.tcl,*.tm') [default: any file]."
	containsDesc = "The comma-separated file names that contain CONTAINS " +
		"(e.g., -c 'readme,README' is the same as " +
		"-g '*readme*,*README*')."
	excludeDesc = "The comma-separated paths to exclude [default: none]."
)
