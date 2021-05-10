[![Test Go status](https://github.com/arduino/libraries-repository-engine/actions/workflows/test-go.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/test-go.yml)
[![Check Go status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go.yml)
[![Check Prettier Formatting status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-prettier-formatting-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-prettier-formatting-task.yml)
[![Spell Check status](https://github.com/arduino/libraries-repository-engine/actions/workflows/spell-check-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/spell-check-task.yml)
[![Check License status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-license.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-license.yml)

## BUILD

```
task go:build
```

## TDD

In order to run the tests, type

```
task go:test
```

## RUN

Create a `config.json` file (or edit example one). Same thing for `repos.txt` file.

Run with `go run sync_libraries.go` or `task go:build` and then `./libraries-repository-engine`

## Security

If you think you found a vulnerability or other security-related bug in this project, please read our
[security policy](https://github.com/arduino/libraries-repository-engine/security/policy) and report the bug to our Security Team 🛡️
Thank you!

e-mail contact: security@arduino.cc
