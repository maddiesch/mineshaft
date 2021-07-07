package main

import (
	"log"
	"os"
)

var (
	// Logger is the shared standard logger for Mineshaft
	Logger = log.New(os.Stderr, "[mineshaft] ", 0)
)
