package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// TODO: remove hard coded paths
const logFileName = "/home/jud/projects/automock/automock.log"

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	Log("UCI Engine Started")
}

func Log(line string) {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	logFile.WriteString(fmt.Sprintf("[%s] %s", time.Now().Format("2006 "+time.StampMilli), line))
}
