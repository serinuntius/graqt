package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	// ADD TWO package
	"github.com/gin-gonic/gin"
	"github.com/serinuntius/graqt"
)

var (
	// ADD TWO VARIABLES
	traceEnabled2 = os.Getenv("GRAQT_TRACE")
	driverName2   = "mysql"

	db2 *sql.DB
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {

	// ADD IF STATEMENT
	if traceEnabled2 == "1" {
		driverName2 = "mysql-tracer"
		// SET LogPath IF YOU'D LIKE TO CHANGE LogPath.
		// DEFAULT IS query.log and request.log
		 graqt.SetQueryLogger("log2/query.log")
		 graqt.SetRequestLogger("log2/request.log")
	}

	var err error

	// CHANGE DRIVER NAME
	db2, err = sql.Open(driverName2, "root:@/graqt")
	if err != nil {
		panic(err)
	}
	defer db2.Close()
	r := gin.Default()

	if traceEnabled2 == "1" {
		r.Use(graqt.RequestIdForGin())
	}

	r.GET("/user", createUser2)
	r.Run(":8080")
}

func createUser2(c *gin.Context) {
	// PLEASE USE CONTEXT
	stmt, _ := db2.PrepareContext(c, "INSERT INTO `user` (email,age) VALUES (?, ?)")
	t1 := time.Now().UnixNano()
	age := rand.Intn(80)

	stmt.ExecContext(c, fmt.Sprintf("hoge%d@hoge.com", t1), age)

	c.String(200, "ok")
}
