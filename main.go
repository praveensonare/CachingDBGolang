package main

import (
	"fmt"
	"strconv"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/go-redis/redis"
)

func main() {
	// search variable
	searchId := 5

	// Open Redis Connection
	redisClient := newRedisClient()
	result, err := redisPing(redisClient)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	// Search value in redis
	nameInRedis, err = redisClient.Get("name_" + strconv.Itoa(searchId)).Result()
	if err != nil { // unexpected error
		fmt.Println(err)
	} else if err == redis.Nil { //key name_[searchId] does not exist in redis
		// Open MySQL Connection
		db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test")
		if err != nil {
			panic(err.Error())
		}
		defer db.Close()

		// Prepare statement for reading data
		stmtOut, err := db.Prepare("SELECT name FROM test WHERE id = ?") // ? is the variable placeholder
		if err != nil {
			panic(err.Error())
		}
		defer stmtOut.Close() // Close the statement when we leave main() or the program terminates

		// Query the rows that has id more than 5
		rows, err := stmtOut.Query(searchId)
		if err != nil {
			panic(err.Error())
		}

		// Process each rows accordingly
		numRows = 0
		for rows.Next() {
			var nameInSQL string
			err = rows.Scan(&nameInSQL)
			if err != nil {
				panic(err.Error())
			}
			fmt.Printf("name is %s\n", nameInSQL)
			numRows = numRows + 1
		}
		if numRows == 0 {
			fmt.Printf("corresponding name is not found\n")
		} else { //key name_[searchId] exists in redis
			fmt.Printf("name is %s\n", nameInRedis)
		}
	}

}

func newRedisClient() * redis.Client{
		redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "", // no password
		DB: 0,        // default DB
	})

	return redisClient
}

func redisPing(client*redis.Client)(string, error) {
	result, err := client.Ping().Result()
	if err != nil {
		return "", err
	} else {
		return result, nil
	}
}
