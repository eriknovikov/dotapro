package main

import "errors"

var ErrNoNewMatches = errors.New("no new matches have been inserted since last scrape")
