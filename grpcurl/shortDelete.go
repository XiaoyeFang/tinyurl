package grpcurl

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"strings"
	"url-shortener/config"
	"url-shortener/protos"
)

func EnerDelete(url, shortener string) (*protos.DelEnerReply, error) {

	reply := &protos.DelEnerReply{}
	if shortener == "" && url == "" {
		reply.Result = "Fill in at least one item"
		return reply, errors.New("The lack of necessary conditions, this will delete the entire table")
	}

	db, err := config.CreateDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	var stmt *sql.Rows
	switch {

	case url == "":
		fmt.Println("DELETE FROM shortener WHERE shortener = $1")
		shortener = strings.Replace(shortener, config.UrlConfig.Options.Prefix, "", -1)
		stmt, err := db.Prepare("DELETE FROM shortener WHERE shortener = $1")
		if err != nil {
			glog.V(0).Infoln(err)
		}

		_, err = stmt.Exec(shortener)

	case shortener == "":

		stmt, err := db.Prepare("DELETE FROM shortener WHERE url =$1")
		if err != nil {
			glog.V(0).Infoln(err)
		}

		_, err = stmt.Exec(url)

	default:
		shortener = strings.Replace(shortener, config.UrlConfig.Options.Prefix, "", -1)
		stmt, err := db.Prepare("DELETE FROM shortener WHERE shortener = $1 AND url =$2")
		if err != nil {
			glog.V(0).Infoln(err)
		}

		_, err = stmt.Exec(shortener, url)

	}
	glog.V(2).Infoln(stmt)

	if err != nil {
		glog.V(0).Infoln(err)
		reply.Result = err.Error()
	} else {
		fmt.Println("delete form postgre success")
		reply.Result = "delete form postgre success"
	}
	return reply, err
}
