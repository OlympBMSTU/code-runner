package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
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
	cfg, _ := config.GetConfigInstance()
	fmt.Println(cfg)
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
			// fmt.Println(i, line)
			record := GetFileRecord(line)
			fileRecords[i-1] = record
		}

	}
	// fmt.Println(fileRecords)

	// for _, fileRecord := range fileRecords {
	// 	fmt.Println(fileRecord.TASKID)
	// }
	return fileRecords, nil
}

func RunR(fName string, input []string) string {
	var cmd *exec.Cmd
	interp := ""
	log := logger.GetLogger()
	if len(interp) > 0 {
		cmd = exec.Command("python", fName)
	} else {
		cmd = exec.Command(fName)
	}

	// buffer := bytes.Buffer{}
	// buffer.Write([]byte("df"))
	// cmd.Start()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Error("Error stdin pipo", err)
		return ""
	}
	defer stdin.Close()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("Error stdout pipo", err)
		return ""
	}
	defer stdout.Close()

	// inStr := ""
	// for _, str := range input {
	// 	inStr += str
	// 	// stdin.Write([]byte(str + "\n"))
	// }
	// _, err := buffer.Write([]byte(inStr + "\n\n\n")) //stdin.Write([]byte(inStr + "\n\n"))

	// cmd.Stdin = &buffer
	// cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		log.Error("Error start", err)
		return ""
	}

	// buffer.Write([]byte("5\n"))
	inStr := ""
	for _, str := range input {
		inStr += str
		// stdin.Write([]byte(str + "\n"))
	}
	n, err := stdin.Write([]byte(inStr + "\n5\n"))
	stdin.Close()
	fmt.Println("N: ", n)
	if err != nil {
		log.Error("Error write stdin", err)
		return ""
	}

	// cmd.ProcessState
	// stdin.Close()
	// stdin.Write([]byte("\n\n"))
	// n, err = stdin.Write([]byte("556\n")

	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Error("Error read stdout", err)
		return ""
	}

	if err = cmd.Wait(); err != nil {
		fmt.Print(err)
		return ""
	}
	fmt.Println(string(b))
	return "" //string(b)
}

func main() {
	cfg, err := config.GetConfigInstance()
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

	for _, fileRecord := range fileRecords {
		if fileRecord.UID == "982" {
			if fileRecord.TASKID == "42" {
				testData := *FindTestData(answers, fileRecord.TASKID)
				for _, d := range testData {
					RunR("/home/imber/file_storage/compiled/982/3_42", d.Input)
				}
			}
		}
	}

	// err = LoopByFiles(fileRecords, answers)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(fileRecords)
	// fmt.Println(answers)
	// fmt.Println(cfg)
}
