package grpcurl

import (
	"database/sql"
	"fmt"
	"strings"
	"url-shortener/config"
	"url-shortener/protos"
)

func EnerQuery(url, shortener string, page, pageSize int32) (*protos.QueryReply, error) {
	reply := &protos.QueryReply{}
	var count int32
	if pageSize <= 0 {
		pageSize = 20
	}
	if page <= 0 {
		page = 1
	}

	db, err := config.CreateDatabase()
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	shortener = strings.Replace(shortener, config.UrlConfig.Options.Prefix, "", -1)
	//判断url和clientType是否为空
	var rows *sql.Rows
	var countrows *sql.Rows
	switch {
	case shortener == "" && url != "":

		rows, err = db.Query(
			"SELECT shortener,clientType,url,uid,createtime FROM shortener WHERE url = $1  ORDER BY createtime ASC limit $2 offset $3",
			url, pageSize, (page-1)*pageSize)

		countrows, err = db.Query(
			"SELECT COUNT(uid) FROM shortener  WHERE url = $1 GROUP BY uid ", url)
		if err != nil {
			fmt.Println("db.Query SELECT", err)
			return nil, err
		}

	case url == "" && shortener != "":

		rows, err = db.Query(
			"SELECT shortener,clientType,url,uid,createtime FROM shortener WHERE shortener =$1  ORDER BY createtime ASC limit $2 offset $3",
			shortener, pageSize, (page-1)*pageSize)
		countrows, err = db.Query(
			"SELECT COUNT(uid) FROM shortener  WHERE shortener = $1 GROUP BY uid ", url)
		if err != nil {
			fmt.Println("db.Query SELECT", err)
			return nil, err
		}
	case url == "" && shortener == "":

		rows, err = db.Query(
			"SELECT shortener,clientType,url,uid,createtime FROM shortener  ORDER BY createtime ASC limit $1 offset $2",
			pageSize, (page-1)*pageSize)

		countrows, err = db.Query(
			"SELECT COUNT(uid) FROM shortener GROUP BY uid ")

		if err != nil {
			fmt.Println("db.Query SELECT", err)
			return nil, err
		}

	default:

		rows, err = db.Query(
			"SELECT shortener,clientType,url,uid,createtime FROM shortener WHERE url = $1 AND shortener =$2 GROUP BY uid limit $3 offset $4",
			url, shortener, pageSize, (page-1)*pageSize)
		countrows, err = db.Query(
			"SELECT COUNT(uid) FROM shortener  WHERE url = $1 AND shortener =$2 GROUP BY uid ", url, shortener)

		if err != nil {
			fmt.Println("db.Query SELECT", err)
			return nil, err
		}
	}

	if rows != nil {

		for rows.Next() {
			re := protos.ShortDetailReply{}

			err = rows.Scan(&re.Shortener, &re.ClientType, &re.Url, &re.Uid, &re.CreateTime)
			if err != nil {
				fmt.Println("rows.Scan", err)
				return reply, err
			}
			//fmt.Println("shortener==", re)
			re.Shortener = config.UrlConfig.Options.Prefix + re.Shortener
			reply.AppInfoList = append(reply.AppInfoList, &re)
		}
	}

	if countrows != nil {

		for countrows.Next() {
			var num int

			err = countrows.Scan(&num)
			if err != nil {
				fmt.Println("rows.Scan", err)
				return reply, err
			}
			count++
		}
	}
	reply.Count = count
	return reply, err
}

func QueryCount() {

	db, err := config.CreateDatabase()
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(
		"SELECT uid, COUNT(uid) FROM shortener WHERE  uid NOT BETWEEN 5 AND 15 GROUP BY uid ",
	)
	if err != nil {
		fmt.Println("db.Query SELECT", err)
		return
	}
	if rows != nil {

		for rows.Next() {
			var uid string
			var count int
			err = rows.Scan(&uid, &count)
			if err != nil {
				fmt.Println("rows.Scan", err)
			}
			fmt.Println("uid", uid, "  count", count)
		}
	}

}
