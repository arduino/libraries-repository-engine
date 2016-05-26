BUILD
----------------------------

```
go get github.com/google/go-github/github
go get github.com/vaughan0/go-ini
go get github.com/blang/semver
go get github.com/stretchr/testify
go get github.com/arduino/arduino-modules/git

go build arduino.cc/repository/libraries-repository-engine
```

You may want to setup git to allow "go get" from private repos, in this case you must
generate a personal access token from github that grants "repo access" permissions and do:

```
git config --global url."https://YOUR_ACCESS_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"
```

the configuration will be saved inside `~/.gitconfig` as:

```
...
[url "https://YOUR_ACCESS_TOKEN:x-oauth-basic@github.com/"]
	insteadOf = https://github.com/
...
```

TDD
----------------------------

In order to run the tests, type

```
go test -v ./src/arduino.cc/repository/libraries/test/...
```

RUN
----------------------------

Create a `config.json` file (or edit example one). Same thing for `repos.txt` file.

Run with `go run sync_libraries.go` or `go build` and then `./libraries-repository-engine`
