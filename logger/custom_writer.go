package logger

import (
	"fmt"
	"os"
)

type LogWriter struct {
	output []*os.File
}

func (wr LogWriter) Write(data string) {
	for _, file := range wr.output {
		fmt.Fprintln(file, data)
		fmt.Fprintln(file, "-------------------------------------------------------------------------")
	}
}

func (wr *LogWriter) AddLogFile(file *os.File) {
	wr.output = append(wr.output, file)
}
