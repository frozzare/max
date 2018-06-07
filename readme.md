# max [![Build Status](https://travis-ci.org/frozzare/max.svg?branch=master)](https://travis-ci.org/frozzare/max) [![Build status](https://ci.appveyor.com/api/projects/status/2m3n5qwam5vv1xvf?svg=true)](https://ci.appveyor.com/project/frozzare/max) [![GoDoc](https://godoc.org/github.com/frozzare/max?status.svg)](http://godoc.org/github.com/frozzare/max) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/max)](https://goreportcard.com/report/github.com/frozzare/max)

Max is a YAML-based task runner.

Check out the [examples](https://github.com/frozzare/max/tree/master/examples).

## Installation

```
go get -u github.com/frozzare/max
```

or using [homebrew](https://brew.sh/).

```
brew install frozzare/tap/max
```

## Usage

Running `max help` will print help output.

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

Default task can be changed by using `MAX_DEFAULT_TASK` environment variable.

```
$ MAX_DEFAULT_TASK=custom max
```

## Task output

Starting and finished logs:

```
$ max hello
Starting task hello
Hello
Finished task hello
```

Minimal logs (quiet flag):

```
$ max hello -q
Hello
```

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
$ max hello -q
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
$ max hello -q
Hello default

$ max hello -q --name max
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
$ max -q
Hello default

$ max default -q --name max
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
$ max hello -q
Hello default

$ max hello -q --name max
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
$ max hello -q
Hello default

$ max hello -q --name max
Hello max
```

## Docker

Tasks can be runned in docker images, you need to configure docker for each task.

Not tested with windows containers (yet, pull request?).

```yaml
tasks:
  build:
    docker:
      image: golang:1.10
      volumes:
        - .:/go/src/app
      working_dir: /go/src/app
    commands:
      - go build -o main
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
  task: task id (os specific tasks can be loaded before real task id, e.g build_windows is loaded when build is called on windows)
    args: Arguments that all tasks can use. Key/Value map that can be used with --key flag.
    deps: [task] # task dependencies, e.g [build, that]
    dir: Custom directory to execute commands in. Default is where the max file is located.
    docker: # docker config
      auth: # private registry auth
        email:
        username:
        password:
      entrypoint: docker entrypoint
      image: docker image
      volumes:
        - single/multi-line array of docker volumes
      working_dir: docker working directory
    interval: task interval (cron format)
    summary: task summary
    tasks:
      - single/multi-line array of tasks to run
    commands:
      - single/multi-line array of commands to run (go text template)
      - access environment variables via $NAME
    status:
      - single/multi-line array of commands to run to test that the task is up to date.
      - (test -e main)
    usage: string of usage text, e.g "[--name]"
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)
