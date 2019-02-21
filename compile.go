package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/OlympBMSTU/code-runner/logger"
)

const OK = 0
const ERROR_COMPILE = 1

type CompileResult struct {
	State     int
	ErrorInfo string
	RunData   RunnerRecord
}

type RunnerRecord struct {
	FilePath      string // full path
	EXID          string
	Interpretator string
}

func Compile(record FileRecord) (*CompileResult, error) {
	log := logger.GetLogger()
	var compiler, outName string
	var res CompileResult
	res.State = OK

	switch record.TYPE {
	// fpc файл -o исполняемый_файл
	case "lpr", "pas":
		compiler = "fpc"
	case "cpp":
		compiler = "g++"
	case "c":
		compiler = "gcc"
	default:
		log.Error("NO SUCH COMPILER "+record.UID+" "+record.TASKID, nil)
		return nil, errors.New("No compiler")
	}

	fName := record.FName
	ext := filepath.Ext(fName)
	outName = record.FName[0 : len(fName)-len(ext)]
	var cmd *exec.Cmd
	if compiler == "fpc" {
		cmd = exec.Command(compiler, fName, "-o"+outName)
	} else {
		cmd = exec.Command(compiler, fName, "-o", outName)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("Error create stdout pipe", err)
		return nil, err
	}
	defer stdout.Close()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Error("Error create stderrpipe", err)
		return nil, err
	}
	defer stderr.Close()

	err = cmd.Start()
	if err != nil {
		log.Error("Error start command", err)
		return nil, err
	}

	stdOutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Error("Error read from stdout", err)
		return nil, err
	}

	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Error("Error read stderr", err)
		return nil, err
	}

	if err = cmd.Wait(); err != nil {
		outStr := string(stdOutBytes)
		errStr := string(stderrBytes)

		if len(outStr) > 0 {
			res.ErrorInfo = outStr
		} else {
			res.ErrorInfo = errStr
		}

		log.Error(res.ErrorInfo, nil)
		res.State = ERROR_COMPILE
	}

	var runRec RunnerRecord
	runRec.FilePath = outName
	runRec.EXID = record.TASKID
	res.RunData = runRec
	return &res, nil
}

func MakeExecutable(record FileRecord) (*RunnerRecord, error) {
	var runRec RunnerRecord
	if record.TYPE != "py" {
		res, err := Compile(record)
		if err != nil {
			return nil, err
		}
		if res.State != OK {
			return nil, errors.New(res.ErrorInfo)
		} else {
			return &res.RunData, nil
		}
	} else {
		runRec.Interpretator = "python"
		runRec.FilePath = record.FName
	}
	return &runRec, nil
}

func LoopByFiles(files []FileRecord, answers []Answer) error {
	// userService := GetService()
	file, err := os.Create("/home/imber/test_result.txt")
	if err != nil {
		panic(err)
	}
	count_to_run := 0
	runnerRecords := make([]RunnerRecord, 0)
	for _, record := range files {
		dbRes := DbTestResult{}
		/////////////////////////////////////////////
		uid, _ := strconv.Atoi(record.UID)
		exID, _ := strconv.Atoi(record.TASKID)
		dbRes.UserID = uid
		dbRes.ExerciseID = exID
		dbRes.FileNameOriginal = record.FName
		//////////////////////////////////////////

		compiled := true
		compilerOutput := ""
		err := MoveFile(&record)
		if err != nil {
			return err
		}

		dbRes.FIleName = record.FName
		runnerRecord, err := MakeExecutable(record)
		if err == nil {
			count_to_run++
			runnerRec := *runnerRecord
			runnerRecords = append(runnerRecords, runnerRec)
			testData := *FindTestData(answers, record.TASKID)

			// testRes := RunTests(runnerRec, testData)
			var testRes TestResult
			if record.UID == "982" {
				fmt.Println(record.FName)
				testRes = RunTests(runnerRec, testData)
				// fmt.Println(res)
				// break

			}
			dbRes.TestResults = testRes.TestResults
			dbRes.TotalMark = testRes.TotalMark
		} else {
			log.Print(err.Error())
			// Write db
			compiled = false
			compilerOutput = err.Error()
		}

		dbRes.Compiled = compiled
		dbRes.CompileOutput = compilerOutput
		bytes, _ := json.Marshal(dbRes)

		if record.UID == "982" { // fmt.Println(string(bytes))
			fmt.Fprint(file, string(bytes))
			fmt.Fprint(file, "\n")
		}
		// userService.Save(dbRes)
		// fmt.Println(runnerRecord)
	}

	fmt.Println("----------------------------------------------")
	fmt.Println("----------------------------------------------")
	fmt.Println("----------------------------------------------")

	return nil
}
