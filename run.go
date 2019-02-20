package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"time"

	"github.com/OlympBMSTU/code-runner/config"
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
	cfg, _ := config.GetConfigInstance()

	fmt.Print(cfg.TimeOutSec)
	log := logger.GetLogger()
	if len(interp) > 0 {
		// cmd = exec.CommandContext(ctx, "python", fName)
		cmd = exec.Command("python", fName)
	} else {
		cmd = exec.Command(fName) // exec.CommandContext(ctx, fName)
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

type TestResult struct {
}

func RunTests(runRec RunnerRecord, testData []AnswerStrcut) TestResult {
	totalMark := 0
	for i, test := range testData {
		res := Run(runRec.FilePath, runRec.Interpretator, test.Input, test.Output)
		if !res.IsError {
			if res.Output == test.Output[i] {
				totalMark += test.Mark
			}
		}
	}

	return TestResult{}
}

// var cmd *exec.Cmd
// cfg, _ := config.GetConfigInstance()

// // ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeOutSec))
// // defer cancel()
// fmt.Print(cfg.TimeOutSec)
// log := logger.GetLogger()
// if len(interp) > 0 {
// 	// cmd = exec.CommandContext(ctx, "python", fName)
// 	cmd = exec.Command("python", fName)
// } else {
// 	cmd = exec.Command(fName) //exec.CommandContext(ctx, fName) //exec.Command(fName)
// }

// // buffer := bytes.Buffer{}
// // buffer.Write([]byte("df"))
// // cmd.Start()
// // context,

// stdin, err := cmd.StdinPipe()
// if err != nil {
// 	log.Error("Error stdin pipo", err)
// 	return ""
// }
// defer stdin.Close()
// stdout, err := cmd.StdoutPipe()
// if err != nil {
// 	log.Error("Error stdout pipo", err)
// 	return ""
// }
// defer stdout.Close()

// // inStr := ""
// // for _, str := range input {
// // 	inStr += str
// // 	// stdin.Write([]byte(str + "\n"))
// // }
// // _, err := buffer.Write([]byte(inStr + "\n\n\n")) //stdin.Write([]byte(inStr + "\n\n"))

// // cmd.Stdin = &buffer
// // cmd.Stdout = os.Stdout
// // cmd.
// err = cmd.Start()
// if err != nil {
// 	log.Error("Error start", err)
// 	return ""
// }

// // buffer.Write([]byte("5\n"))
// inStr := ""
// for _, str := range input {
// 	inStr += str
// 	// stdin.Write([]byte(str + "\n"))
// }
// n, err := stdin.Write([]byte(inStr + "\n5\n"))
// stdin.Close()
// fmt.Println("N: ", n)
// if err != nil {
// 	log.Error("Error write stdin", err)
// 	return ""
// }

// // stdin.Close()
// // stdin.Write([]byte("\n\n"))
// // n, err = stdin.Write([]byte("556\n")

// b, err := ioutil.ReadAll(stdout)
// fmt.Println("Process state: ", cmd.ProcessState)
// if err != nil {
// 	log.Error("Error read stdout", err)
// 	return ""
// }

// done := make(chan error, 1)
// go func() {
// 	done <- cmd.Wait()
// }()
// select {
// case <-time.After(3 * time.Second):
// 	if err := cmd.Process.Kill(); err != nil {
// 		log.Error("failed to kill process: ", err)
// 	}
// 	log.Info("process killed as timeout reached")
// case err := <-done:
// 	if err != nil {
// 		log.Error("process finished with error = %v", err)
// 	}
// 	log.Info("process finished successfully")
// }

// // if ctx.Err() == context.DeadlineExceeded {
// // 	fmt.Println("Command timed out")
// // 	return ""
// // }

// // if err = cmd.Wait(); err != nil {
// // 	fmt.Print(err)
// // 	return ""
// // }
// fmt.Println(string(b))
// return "" //string(b)
