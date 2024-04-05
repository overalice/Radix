package radix

import (
	"log"
	"os"
)

var (
	Info  = log.New(os.Stdout, "\033[34m[info   ]\033[0m ", log.LstdFlags).Printf
	Warn  = log.New(os.Stdout, "\033[33m[warning]\033[0m ", log.LstdFlags|log.Lshortfile).Printf
	Error = log.New(os.Stdout, "\033[31m[error  ]\033[0m ", log.LstdFlags|log.Lshortfile).Printf
)
