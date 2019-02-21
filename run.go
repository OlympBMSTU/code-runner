package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/OlympBMSTU/code-runner/logger"
)

type RunResult struct {
	Output      string
	IsError     bool
	ErrorOutput string
	Error       error
}

func Run(fName string, interp string, input []string, answer []string) RunResult {
	var cmd *exec.Cmd
	// cfg, _ := config.GetConfigInstance()

	path := filepath.Dir(fName)
	id := filepath.Base(fName)
	testFile := path + "/" + "input" + id

	inputFile, _ := os.Create(testFile)

	defer inputFile.Close()
	for _, line := range input {
		fmt.Fprint(inputFile, line)
	}

	log := logger.GetLogger()
	if len(interp) > 0 {
		cmd = exec.Command("python", fName) //, "<", testFile)
	} else {
		cmd = exec.Command(fName) //, "< ", testFile) // exec.CommandContext(ctx, fName)
	}

	done := make(chan RunResult, 1)

	go func() {
		var runRes RunResult

		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Error("Error stdin pipo", err)
			done <- runRes
			return
		}
		defer stdin.Close()
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Error("Error stdout pipo", err)
			done <- runRes
			return
		}
		defer stdout.Close()

		err = cmd.Start()
		if err != nil {
			log.Error("Error start", err)
			done <- runRes
			return
		}

		// todo refactor
		inStr := ""
		for _, str := range input {
			inStr += str
		}

		// we send new data for prevent stdin
		_, err = stdin.Write([]byte(inStr + "\n5\n"))
		if err != nil {
			log.Error("Error write stdin", err)
			done <- runRes
			return
		}

		b, err := ioutil.ReadAll(stdout)
		if err != nil {
			(&runRes).Error = err
			(&runRes).IsError = true
			log.Error("Error read stdout", err)
			done <- runRes
			return
		}

		if err = cmd.Wait(); err != nil {
			fmt.Print(err)
			(&runRes).Error = err
			(&runRes).IsError = true
			(&runRes).ErrorOutput = err.Error()
			done <- runRes
			return
		}

		(&runRes).IsError = false
		(&runRes).Output = string(b)

		done <- runRes
	}()

	select {
	case <-time.After(15 * time.Second):
		res := RunResult{
			IsError:     true,
			ErrorOutput: "Max Time limit, or infinite loop",
		}
		if err := cmd.Process.Kill(); err != nil {
			log.Error("failed to kill process: ", err)
		}

		log.Info("process killed as timeout reached")
		return res
	case runRes := <-done:
		if runRes.IsError {
			log.Error("process finished with error = %v", runRes.Error)
		} else {
			log.Info("OK: " + runRes.Output)
		}
		log.Info("process finished successfully")
		return runRes
	}
}

type UserTest struct {
	ID     int
	Input  []string
	Output []string
	Mark   int
}

type TestResult struct {
	TotalMark   int
	TestResults []UserTest
}

func RunTests(runRec RunnerRecord, testData []AnswerStrcut) TestResult {
	totalMark := 0
	testRes := make([]UserTest, len(testData))
	for i, test := range testData {
		testRes[i].ID = test.ID
		testRes[i].Input = test.Input
		testRes[i].Mark = 0
		fmt.Println(fmt.Sprintf("Running exercise: %s And Test %d\n Path %s\n Its input: %v\n Expectiong output: %v ", runRec.EXID, test.ID, runRec.FilePath, test.Input, test.Output))
		res := Run(runRec.FilePath, runRec.Interpretator, test.Input, test.Output)
		var output string
		if !res.IsError {
			resData := strings.Trim(res.Output, "\n")
			for _, data := range test.Output {
				if data == resData {
					totalMark += test.Mark
					testRes[i].Mark = test.Mark
				}
			}
			output = resData
			fmt.Println(fmt.Sprintf("State OK\n Result %s %d", res.Output, totalMark))
		} else {
			fmt.Println(fmt.Sprintf("State Error\n Result %s", res.ErrorOutput))
			output = res.ErrorOutput
		}
		testRes[i].Output = []string{output}
	}

	return TestResult{
		TotalMark:   totalMark,
		TestResults: testRes,
	}
}
