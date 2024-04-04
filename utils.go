package radix

import (
	"log"
	"os"
)

var (
	info  = log.New(os.Stdout, "\033[34m[info   ]\033[0m ", log.LstdFlags).Printf
	warn  = log.New(os.Stdout, "\033[33m[warning]\033[0m ", log.LstdFlags|log.Lshortfile).Printf
	fault = log.New(os.Stdout, "\033[31m[error  ]\033[0m ", log.LstdFlags|log.Lshortfile).Printf
)
