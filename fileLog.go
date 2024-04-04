package radix

import (
	"log"
	"os"
	"path/filepath"
)

var (
	stlog   *log.Logger
	logPath = "/log"
)

type fileLog string

func init() {
	var err error
	currentDir, err = os.Getwd()
	if err != nil {
		fault(err.Error())
		return
	}

	_, err = os.Stat(filepath.Join(currentDir, logPath))
	if os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Join(currentDir, logPath), 0755)
		if err != nil {
			fault(err.Error())
			return
		}
		info("Created a folder: %s", logPath)
	}
}

func StartFileLog(destination string) {
	stlog = log.New(fileLog(filepath.Join(logPath, destination)), "radix: ", log.LstdFlags)
}

func (fl fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func WriteLog(message string) {
	if stlog == nil {
		fault("Filelog is not enabled")
	} else {
		stlog.Printf("%v\n", message)
	}
}
