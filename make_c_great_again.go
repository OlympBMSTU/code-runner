package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/OlympBMSTU/code-runner/config"
	"github.com/OlympBMSTU/code-runner/logger"
)

const templateHeaderCpp = `#include <iostream>\n
#include <algorithms>\n
`

const templateHeaderC = `#include <stdio.h>\n`
const templateMain = `int main(){\n`

const templateEndMain = `return 0;\n}`

func Exist(data string, arr []string) bool {
	for _, d := range arr {
		if d == data {
			return true
		}
	}
	return false
}

func FindDependecies() []string {
	cfg, err := config.GetConfigInstance()
	err = logger.InitLogger(*cfg)
	if err != nil {
		log.Println(err.Error())
		return []string{}
	}
	fileRecords, err := FilesLoop()
	if err != nil {
		return []string{}
	}

	arrayDep := make([]string, 0)
	for _, rec := range fileRecords {
		if rec.TYPE == "c" || rec.TYPE == "cpp" {
			bytes, _ := ioutil.ReadFile(cfg.FilesPath + "/" + rec.FName)
			stringss := strings.Split(string(bytes), "\n")
			for _, str := range stringss {
				if strings.Contains(str, "#include") {
					// fmt.Print(rec.FName)
					lib := str
					// lib := strings.Split(str, " ")[1]
					if !Exist(lib, arrayDep) {
						arrayDep = append(arrayDep, lib)
					}
				}
			}
		}
	}

	for _, lib := range arrayDep {
		fmt.Println(lib)
	}
	return arrayDep
	// fmt.Print(string(b))
}
