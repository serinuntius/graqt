graqt ~go-request-query-logger~
===

`graqt` is middleware to easily record access log and sql log.

## Description
`graqt` record a request_time, request_path, method, and unique request id in access log.

Also `graqt` record a sql with the unique id in query log.

Because Publishing a unique id to check how many times sql was called for each request.


## Demo
query.log
```json
{"level":"info","ts":1525082389.1336403,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO `user` (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525082389130688404@hoge.com"},{"Name":"","Ordinal":2,"Value":57}],"time":0.001157627,"request_id":"5914e629-6746-42c5-b342-702f224f48e1"}
{"level":"info","ts":1525082389.215664,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO `user` (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525082389213632765@hoge.com"},{"Name":"","Ordinal":2,"Value":37}],"time":0.002013919,"request_id":"4e743b2a-59be-4692-b36b-3b6a9ba78b01"}
{"level":"info","ts":1525082389.2982852,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO `user` (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525082389287262677@hoge.com"},{"Name":"","Ordinal":2,"Value":74}],"time":0.011005944,"request_id":"f1d68fb3-5450-48f3-9f46-1ea5ad2370fb"}
{"level":"info","ts":1525082389.37562,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO `user` (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525082389372293507@hoge.com"},{"Name":"","Ordinal":2,"Value":64}],"time":0.003297617,"request_id":"e727db2b-df12-410c-91a4-4dd3117b9452"}
{"level":"info","ts":1525082389.474928,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO `user` (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525082389470871128@hoge.com"},{"Name":"","Ordinal":2,"Value":73}],"time":0.001937477,"request_id":"57df4708-9705-4401-98b2-fad4bacfe585"}
{"level":"info","ts":1525082389.5541484,"caller":"graqt/tracer.go:20","msg":"Exec","query":"INSERT INTO `user` (email,age) VALUES (?, ?)","args":[{"Name":"","Ordinal":1,"Value":"hoge1525082389550824448@hoge.com"},{"Name":"","Ordinal":2,"Value":39}],"time":0.001194376,"request_id":"58957516-c7b2-4f64-8889-b39496e6597f"}
```

access.log
```json
{"level":"info","ts":1525082389.1337767,"caller":"graqt/middleware.go:25","msg":"","time":0.011566198,"id":"5914e629-6746-42c5-b342-702f224f48e1","path":"/user"}
{"level":"info","ts":1525082389.2157433,"caller":"graqt/middleware.go:25","msg":"","time":0.003945689,"id":"4e743b2a-59be-4692-b36b-3b6a9ba78b01","path":"/user"}
{"level":"info","ts":1525082389.2983494,"caller":"graqt/middleware.go:25","msg":"","time":0.014426681,"id":"f1d68fb3-5450-48f3-9f46-1ea5ad2370fb","path":"/user"}
{"level":"info","ts":1525082389.3756797,"caller":"graqt/middleware.go:25","msg":"","time":0.004002398,"id":"e727db2b-df12-410c-91a4-4dd3117b9452","path":"/user"}
{"level":"info","ts":1525082389.47504,"caller":"graqt/middleware.go:25","msg":"","time":0.005889232,"id":"57df4708-9705-4401-98b2-fad4bacfe585","path":"/user"}
{"level":"info","ts":1525082389.554212,"caller":"graqt/middleware.go:25","msg":"","time":0.016009073,"id":"58957516-c7b2-4f64-8889-b39496e6597f","path":"/user"}
```


## Requirement
- dep

## Usage

Add setting to use graqt.

Example
```go
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
	age := rand.Intn(80)

	stmt.ExecContext(ctx, fmt.Sprintf("hoge%d@hoge.com", t1), age)

	w.Write([]byte(`ok`))
}

```

## Install

- add `import "github.com/serinuntius/graqt"`

```bash
$ dep ensure
```

## Contribution
1. Fork ([https://github.com/serinuntius/graqt/fork](https://github.com/serinuntius/graqt/fork))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Run test suite with the `go test ./...` command and confirm that it passes
6. Run `gofmt -s`
7. Create new Pull Request


## Licence

[MIT](https://github.com/serinuntius/graqt/blob/master/LICENCE)
