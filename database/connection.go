package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/pranotobudi/graphql-checkout/config"
)

type PostgresDB struct {
	DB *sql.DB
}

var PostgresDBInstance *PostgresDB
var once sync.Once

func GetDB() *PostgresDB {
	once.Do(func() {
		dbConfig := config.DbConfig()
		psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
			dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.DbName, dbConfig.SSLMode)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			log.Println("Fail to setup Postgres connection!")
			return
			// return nil
			// panic(err)
		}
		// retry connection
		retryCount := 30
		for {
			err := db.Ping()
			if err != nil {
				if retryCount == 0 {
					log.Fatalf("Not able to establish connection to database")
				}

				log.Println(fmt.Sprintf("Could not connect to database. Wait 5 seconds. %d retries left...", retryCount))
				retryCount--
				time.Sleep(5 * time.Second)
			} else {
				break
			}
		}

		PostgresDBInstance = &PostgresDB{DB: db}
	})

	log.Println("Successfully connected to Postgres!")
	return PostgresDBInstance
}

func (pg *PostgresDB) MigrateDB(sqlFile string) error {
	file, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		log.Println("Failed to read file")
		return err
	}
	log.Println("Success to read file")
	ctx, cancelfunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelfunc()

	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		log.Println("query: ", q)
		if q == "" {
			continue
		}
		_, err := pg.DB.ExecContext(ctx, q)
		if err != nil {
			log.Println("SQL statement execution failed: ", err)
			return err
		}
		log.Println("table created")

	}
	log.Println("Success to execute database query")
	return nil
}
