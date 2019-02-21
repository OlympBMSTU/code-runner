package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/OlympBMSTU/code-runner/config"
	"github.com/OlympBMSTU/code-runner/logger"
	"github.com/jackc/pgx"
)

const INSERT_RESULT = "INSERT INTO test_result($1, $2, $3, $4, $5, $6, $7, $8, $9)"

type DbTestResult struct {
	Compiled         bool // what with python
	CompileOutput    string
	TotalMark        int
	UserID           int
	ExerciseID       int
	TestResults      []UserTest
	FIleName         string
	FileNameOriginal string
}

type UserAnswersService struct {
	db *pgx.ConnPool
}

var service *UserAnswersService

func InitUserService(cfg config.Config) {
	// if err != nil {
	// 	log.Print(err)
	// 	return nil, err
	// }
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			User:     cfg.DBUser,
			Port:     uint16(5432),
			Password: cfg.DBPassword,
			Database: cfg.DBName,
		},
		MaxConnections: 5,
	}

	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Println(err)
		panic(err)
	}
	service = &UserAnswersService{
		db: pool,
	}
}

func GetService() *UserAnswersService {
	return service
}

// create table test_results(
//     id serial,
//     u_id int,
//     ex_id int,
//     mark int,
//     error citext
//     compiled bool,
//     compile_output citext,
//     file_name citext,
//     file_name_old citext,
//     run_output jsonb
// )

func (serv UserAnswersService) Save(testResult DbTestResult) error {

	log := logger.GetLogger()
	testResults, err := json.Marshal(testResult.TestResults)
	if err != nil {
		log.Error("Error marshal user answers", err)
		return err
	}
	fmt.Println(testResults)
	return nil
}
