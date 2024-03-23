package radix

import "fmt"

func info(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}

func warn(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}

func error(format string, values ...interface{}) {
	fmt.Printf(format+"\n", values...)
}
