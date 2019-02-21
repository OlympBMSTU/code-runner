package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
)

type AnswerStrcut struct {
	ID     int      `json:"id"`
	Input  []string `json:"input"`
	Output []string `json:"output"`
	Mark   int      `json:"mark"`
}

type Answer struct {
	TASKID  string
	Answer  string
	Answers []AnswerStrcut
}

func FindTestData(answers []Answer, exID string) *[]AnswerStrcut {
	for _, answer := range answers {
		if answer.TASKID == exID {
			return &answer.Answers
		}
	}
	return nil
}

func GetAnswers() ([]Answer, error) {
	bytes, err := ioutil.ReadFile("/home/imber/answers.txt")
	var answers []Answer
	if err != nil {
		log.Print(err.Error())
		return answers, err
	}

	lines := strings.Split(string(bytes), "\n")

	answers = make([]Answer, len(lines))
	for i, line := range lines {
		if i > 1 {
			data := strings.Split(line, "\t")
			var answer Answer
			answer.TASKID = data[0]
			answer.Answer = data[11]
			var answersStr []AnswerStrcut
			err := json.Unmarshal([]byte(data[11]), &answersStr)

			if err != nil {
				log.Print(err.Error())
				return answers, err
			}
			answer.Answers = answersStr
			answers[i] = answer
		}
	}
	return answers, nil
}
