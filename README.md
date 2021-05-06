[![Check Go status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go.yml)

BUILD
----------------------------

```
task go:build
```

TDD
----------------------------

In order to run the tests, type

```
task go:test
```

RUN
----------------------------

Create a `config.json` file (or edit example one). Same thing for `repos.txt` file.

Run with `go run sync_libraries.go` or `task go:build` and then `./libraries-repository-engine`
