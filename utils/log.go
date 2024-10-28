package utils

import (
	"fmt"
	"os"
	"strings"
	"time"
)

var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("/home/jud/projects/stockhuman/stockhuman.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	Log(fmt.Sprintf("%s %s started", "Stockhuman", "1.0"))
}

func Log(line string) {
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	logFile.WriteString(fmt.Sprintf("[%s] %s", time.Now().Format("2006 "+time.StampMilli), line))
}
