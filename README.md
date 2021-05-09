[![Check Go status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go.yml)
[![Check Prettier Formatting status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-prettier-formatting-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-prettier-formatting-task.yml)
[![Spell Check status](https://github.com/arduino/libraries-repository-engine/actions/workflows/spell-check-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/spell-check-task.yml)

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
