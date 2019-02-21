package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/OlympBMSTU/code-runner/config"
	"github.com/OlympBMSTU/code-runner/logger"
)

type FileRecord struct {
	FName  string
	TYPE   string
	UID    string
	PID    string
	VAR    string
	TASKNO string
	TASKID string
}

func GetFileRecord(line string) FileRecord {
	data := strings.Split(line, "\t")
	fmt.Println(data)
	return FileRecord{
		data[0],
		data[1],
		data[2],
		data[3],
		data[4],
		data[5],
		data[6],
	}
}

func FilesLoop() ([]FileRecord, error) {
	// cfg, _ := config.GetConfigInstance()
	// fmt.Println(cfg)
	var fileRecords []FileRecord
	bytes, err := ioutil.ReadFile("/home/imber/result.txt") //cfg.FilesPath + "/result.txt")
	if err != nil {
		log.Println(err.Error())
		return fileRecords, err
	}

	lines := strings.Split(string(bytes), "\n")
	fileRecords = make([]FileRecord, len(lines))
	for i, line := range lines {
		if i > 0 {
			record := GetFileRecord(line)
			fileRecords[i-1] = record
		}

	}
	return fileRecords, nil
}

func main() {
	cfg, err := config.GetConfigInstance()
	InitUserService(*cfg)
	err = logger.InitLogger(*cfg)
	if err != nil {
		log.Println(err.Error())
		return
	}

	fileRecords, err := FilesLoop()
	if err != nil {
		return
	}

	// py := 0
	// c := 0
	// cpp := 0
	// pas := 0

	// for _, rec := range fileRecords {
	// 	switch rec.TYPE {
	// 	case "py":
	// 		py++
	// 	case "c":
	// 		c++
	// 	case "cpp":
	// 		cpp++
	// 	case "pas", "lpr":
	// 		pas++
	// 	}
	// }

	// fmt.Print(py, c, cpp, pas)
	answers, err := GetAnswers()
	if err != nil {
		return
	}

	err = LoopByFiles(fileRecords, answers)
	if err != nil {
		fmt.Println(err.Error())
	}

}
