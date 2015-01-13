Requirement: Install git2go.
----------------------------

Instructions here: https://github.com/libgit2/git2go#installing

Quick recap:

    sudo apt-get install cmake
    go get -d github.com/libgit2/git2go
    cd src/github.com/libgit2/git2go
    git submodule update --init # get libgit2
    make install
    cd -

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
