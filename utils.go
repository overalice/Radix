package radix

import "fmt"

func info(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}

func warn(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}

func fault(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}
