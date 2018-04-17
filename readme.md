# max [![Build Status](https://travis-ci.org/frozzare/max.svg?branch=master)](https://travis-ci.org/frozzare/max) [![GoDoc](https://godoc.org/github.com/frozzare/max?status.svg)](http://godoc.org/github.com/frozzare/max) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/max)](https://goreportcard.com/report/github.com/frozzare/max)

> Work in progress, not production ready

Max is a YAML-based task runner.

## Installation

```
go get -u github.com/frozzare/max
```

## Usage

Running `max help` will show a output of options.

```
Options:
  -a, --all             runs all tasks
  -c, --config string   sets the config file
  -l, --list            lists tasks with summary description
  -o, --once            runs tasks once and ignore interval
  -v, --version         show max version
```

## Task help

```
$ max help [task]
Usage:
  max hello [name]

Summary:
  Hello task
```

## Configuration

The default task is `default`

### Basic task

Config

```yaml
tasks:
  hello:
    summary: Hello task
    commands:
      - echo Hello
```

Output

```
$ max hello
Hello
```

### Task with arguments

Config

```yaml
tasks:
  hello:
    args:
      name: default
    summary: Hello task
    commands:
      - echo Hello {{ .name }}
    usage: "[--name]"
```

Output

```
$ max hello
Hello default

$ max hello --name max
Hello max
```

### Task running other tasks

Config

```yaml
tasks:
  hello:
    args:
      name: default
    summary: Hello task
    commands:
      - echo Hello {{ .name }}
  default:
    tasks:
      - hello
```

Output

```
$ max
Hello default

$ max default --name max
Hello max
```

### Task with global arguments

Config

```yaml
args:
  name: default

tasks:
  hello:
    summary: Hello task
    commands:
      - echo Hello {{ .name }}
```

Output

```
$ max hello
Hello default

$ max hello --name max
Hello max
```

### Include task from other files.

Config `max.yml`

```yaml
tasks:
  hello: !include hello.yml
```

Config `hello.yml`

```yaml
args:
  name: default
summary: Hello task
commands:
  - echo Hello {{ .name }}
```

Output

```
$ max hello
Hello default

$ max hello --name max
Hello max
```

## Max file spec

The default file name is `max.yml` but you can specific another file by using the `--config` flag.

Other supported default files are:

- `max_windows.yml`
- `max_linux.yml`
- `max_darwin.yml`

```yaml
args: Global arguments that all tasks can use. Key/Value map that can be used with --key flag.
tasks:
  task: task id
    args: Arguments that all tasks can use. Key/Value map that can be used with --key flag.
    deps: [task] # task dependencies, e.g [build, that]
    dir: Custom directory to execute commands in. Default is where the max file is located.
    interval: task interval (cron format) (not required)
    summary: task summary (not required)
    tasks:
      - multi-line array of tasks to run
      - (not required)
    commands:
      - multi-line array of commands to run (go string formatting supported with arguments)
      - access environment variables via $NAME
      - (not required)
    usage: string of usage text, e.g "[--name]" (not required)
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)
