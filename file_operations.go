package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/OlympBMSTU/code-runner/config"
)

func CheckDirExist(dirPath string) bool {
	if _, err := os.Stat(dirPath); err != nil {
		return false
	}
	return true
}

func MoveFile(record *FileRecord) error {
	cfg, _ := config.GetConfigInstance()
	file, err := os.Open(cfg.FilesPath + "/" + record.FName)
	if err != nil {
		log.Print(err.Error())
		return err
	}
	defer file.Close()

	fName := record.FName
	ext := filepath.Ext(fName)

	dirName := cfg.ComiledPath + "/" + fName[0:len(fName)-len(ext)]
	if !CheckDirExist(dirName) {
		if err := os.Mkdir(dirName, 0777); err != nil {
			log.Print(err.Error())
			return err
		}
	}

	newName := dirName + "/" + record.TASKNO + "_" + record.TASKID + "." + record.TYPE
	fileNew, err := os.Create(newName)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	defer fileNew.Close()
	_, err = io.Copy(fileNew, file)
	record.FName = newName

	return nil
}
