# libraries-repository-engine

[![Test Go status](https://github.com/arduino/libraries-repository-engine/actions/workflows/test-go-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/test-go-task.yml)
[![Codecov](https://codecov.io/gh/arduino/libraries-repository-engine/branch/main/graph/badge.svg)](https://codecov.io/gh/arduino/libraries-repository-engine)
[![Integration Test status](https://github.com/arduino/libraries-repository-engine/actions/workflows/test-go-integration-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/test-go-integration-task.yml)
[![Check Go status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go-task.yml)
[![Check Python status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-python-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-python-task.yml)
[![Check General Formatting status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-general-formatting-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-general-formatting-task.yml)
[![Check Prettier Formatting status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-prettier-formatting-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-prettier-formatting-task.yml)
[![Check Taskfiles status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-taskfiles.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-taskfiles.yml)
[![Spell Check status](https://github.com/arduino/libraries-repository-engine/actions/workflows/spell-check-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/spell-check-task.yml)
[![Check Markdown status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-markdown-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-markdown-task.yml)
[![Check License status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-license.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-license.yml)
[![Check Go Dependencies status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go-dependencies-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-go-dependencies-task.yml)
[![Check YAML status](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-yaml-task.yml/badge.svg)](https://github.com/arduino/libraries-repository-engine/actions/workflows/check-yaml-task.yml)

This is the tool that generates [the Arduino Library Manager index](http://downloads.arduino.cc/libraries/library_index.json).

Every hour, the automated Library Manager indexer system runs this tool, which:

1. checks every repository in the [Library Manager list](https://github.com/arduino/library-registry) for new [tags](https://git-scm.com/book/en/v2/Git-Basics-Tagging)
1. checks whether those tags meet [the requirements for addition to the index](https://github.com/arduino/library-registry/blob/main/FAQ.md#what-are-the-requirements-for-publishing-new-releases-of-libraries-already-in-the-library-manager-list), publishing [logs](https://github.com/arduino/library-registry/blob/main/FAQ.md#can-i-check-on-library-releases-being-added-to-library-manager)
1. adds entries to the index for compliant tags
1. pushes the updated index to Arduino's download server

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

Run the following command to list the available command line interfaces:

```
./libraries-repository-engine help
```

## Security

If you think you found a vulnerability or other security-related bug in this project, please read our
[security policy](https://github.com/arduino/libraries-repository-engine/security/policy) and report the bug to our Security Team 🛡️
Thank you!

e-mail contact: security@arduino.cc
