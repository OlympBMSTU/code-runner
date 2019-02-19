package main

import (
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/OlympBMSTU/code-runner/logger"
)

type RunResult struct {
	State     string
	Mark      int
	ErrorInfo string
}

func RunCmd(cmd *exec.Cmd) {
	cmd.Start()

}

func Run(fName string, interp string, input []string, answer []string) string {
	var cmd *exec.Cmd
	log := logger.GetLogger()
	if len(interp) > 0 {
		cmd = exec.Command("python", fName)
	} else {
		cmd = exec.Command(fName)
	}

	// buffer := bytes.Buffer{}
	// buffer.Write([]byte("df"))
	// cmd.Start()
	// context,

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
	// cmd.
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

	// stdin.Close()
	// stdin.Write([]byte("\n\n"))
	// n, err = stdin.Write([]byte("556\n")

	b, err := ioutil.ReadAll(stdout)
	fmt.Println("Process state: ", cmd.ProcessState)
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

// 10 50 11 00 3 20

func RunTests(runRec RunnerRecord, testData []AnswerStrcut) RunResult {
	for _, test := range testData {
		fmt.Print(test)

		Run(runRec.FilePath, runRec.Interpretator, test.Input, test.Output)
		// if test succ -> mark += test.Mark

	}

	// fill res
	return RunResult{}
}
