BUILD
----------------------------

go get github.com/google/go-github/github
go get github.com/robfig/cron
go get github.com/vaughan0/go-ini
go get github.com/blang/semver
go get github.com/stretchr/testify

go build -o libraries-repository-engine

TDD
----------------------------

In order to run the test, run

```
go test -v ./test/...
```

RUN
----------------------------

Create a `config.json` file (or edit example one). Same thing for `repos.txt` file.

Run with `go run sync_libraries.go` or `go build` and then `./libraries-repository-engine`
