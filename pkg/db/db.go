// Пакет для работы с БД приложения GoNews.
package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	username = "root"
	password = "root"
	hostname = "127.0.0.1"
	port     = 3306
	dbName   = "newsdb"
)

// База данных.
type DB struct {
	pool *sql.DB
}

type PostComment struct {
	NewsId    int
	Comment   string
	IsComment int
}

type ValidateCommet struct {
	Id      int    `json:"id"`
	Comment string `json:"comment"`
}

type row struct {
	ID int64 `field:"id"`
}

func New() (*DB, error) {

	connString := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		username,
		password,
		hostname,
		port,
		dbName,
	)

	log.Println("connString: ", connString)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, fmt.Errorf("mysql err: %s", err)
	}

	return &DB{
		pool: db,
	}, nil
}

func (db *DB) AddComment(comment PostComment) error {
	query := fmt.Sprintf("INSERT INTO comment (`comment`, `forId`, `for_comment`) VALUES ('%s', %d, %d)",
		comment.Comment,
		comment.NewsId,
		comment.IsComment,
	)

	r, err := db.pool.Exec(query)
	if err != nil {

		return err
	}

	id, err := r.LastInsertId()
	if err != nil {

		return err
	}

	url := "http://localhost:83/validate"

	vc := ValidateCommet{
		Id:      int(id),
		Comment: comment.Comment,
	}

	jsonStr, err := json.Marshal(vc)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		return err
	}

	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return nil
}
