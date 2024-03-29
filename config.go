package radix

import (
	"bufio"
	"os"
	"strings"
)

var config map[string]string

func init() {
	config = make(map[string]string)

	file, err := os.Open("config.txt")
	if err != nil {
		fault("Failed to open config file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			warn("Failed to parse %s", line)
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		config[key] = value
	}
	if err := scanner.Err(); err != nil {
		fault(err.Error())
	}
}
