package postgres

import (
	"database/sql"
	"fmt"

	// This loads the postgres drivers.
	_ "github.com/lib/pq"

	"time"
	"url-shortener/base62"
	"url-shortener/storage"
)

// New returns a postgres backed storage service.
func New(host, port, user, password, dbName string) (storage.Service, error) {
	// Connect postgres
	connect := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", connect)
	if err != nil {
		return nil, err
	}

	// Ping to connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Create table if not exists
	strQuery := "CREATE TABLE IF NOT EXISTS shortener (uid serial NOT NULL, url VARCHAR not NULL, " +
		"visited boolean DEFAULT FALSE, shortener VARCHAR not NULL,clienttype VARCHAR not NULL, " +
		"createtime VARCHAR not NULL,count INTEGER DEFAULT 0);"

	_, err = db.Exec(strQuery)
	if err != nil {
		return nil, err
	}
	return &Postgres{db}, nil
}

type Postgres struct{ db *sql.DB }

func (self *Postgres) Save(url, clientType string) (string, error) {
	var shortener string
	//查重操作
	rows, err := self.db.Query("SELECT shortener FROM shortener WHERE url = $1 AND clienttype =$2", url, clientType)
	if err != nil {
		fmt.Println("p.db.Query", err)
	}
	if rows != nil {

		for rows.Next() {

			err = rows.Scan(&shortener)
			if err != nil {
				fmt.Println("rows.Scan", err)
				return "", err
			}
		}
		if shortener == "" {
			fmt.Println("数据库中无数据")
			shortener, err = base62.Transform(time.Now().String() + url)
			fmt.Println("shortener", shortener)
			var id int
			fmt.Println(clientType)
			err = self.db.QueryRow("INSERT INTO shortener(url,visited,shortener,clienttype,createtime,count) VALUES($1,$2,$3,$4,$5,$6) returning uid;", url, false, shortener, clientType, time.Now().String(), 0).Scan(&id)
			if err != nil {
				fmt.Println("p.db.QueryRow INSERT", err)
				return "", err
			}

		}
	}
	return shortener, nil
}

/*
If you see here that you are very amazing, because I also don't understand the following code
*/
func (self *Postgres) Load(code string) (string, error) {

	var url string
	err := self.db.QueryRow("update shortener set visited=true, count = count + 1 where shortener=$1 RETURNING url", code).Scan(&url)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (self *Postgres) LoadInfo(code string) (*storage.Item, error) {

	fmt.Println("LoadInfo code=", code)

	var item storage.Item
	err := self.db.QueryRow("SELECT url, visited, count FROM shortener where shortener=$1 limit 1", code).
		Scan(&item.URL, &item.Visited, &item.Count)
	if err != nil {
		return nil, err
	}
	fmt.Println("item.URL", item.URL)
	return &item, nil
}

func (self *Postgres) Close() error { return self.db.Close() }
