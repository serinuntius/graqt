package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"

	// ADD TWO package
	"github.com/justinas/alice"
	"github.com/serinuntius/graqt"
)

var (
	// ADD TWO VARIABLES
	traceEnabled = true
	driverName   = "mysql"

	db *sql.DB
)

func init() {
	rand.Seed(time.Now().UnixNano())
}
func main() {

	// ADD IF STATEMENT
	if traceEnabled {
		driverName = "mysql-tracer"
		// SET LogPath IF YOU'D LIKE TO CHANGE LogPath.
		// DEFAULT IS query.log and request.log
		// graqt.SetQueryLogger("log/query.log")
		// graqt.SetRequestLogger("log/request.log")
	}

	var err error

	// CHANGE DRIVER NAME
	db, err = sql.Open(driverName, "root:@/graqt")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/user", createUser)

	// ADD MIDDLEWARE
	var chain alice.Chain
	if traceEnabled {
		chain = alice.New(graqt.RequestId)
	} else {
		chain = alice.New()
	}

	// CHANGE THIS LINE
	http.ListenAndServe(":8080", chain.Then(mux))
}

func createUser(w http.ResponseWriter, r *http.Request) {
	// PLEASE USE CONTEXT
	ctx := r.Context()
	stmt, _ := db.PrepareContext(ctx, "INSERT INTO `user` (email,age) VALUES (?, ?)")
	t1 := time.Now().UnixNano()
	age := rand.Int()

	stmt.ExecContext(ctx, fmt.Sprintf("hoge%d@hoge.com", t1), age)

	w.Write([]byte(`ok`))
}
