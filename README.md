# kasher

**kasher** is a Go-based CLI tool for defining, running, and caching named shell tasks. It helps you wrap shell commands as reusable tasks, cache their output for a configurable expiration time, and manage them interactively.

## Features

- Define named tasks that wrap shell commands
- Cache task output for a set expiration time
- Store configuration (location based on OS and user setup, MacOS: ~/Library/Application Support/kasher/config.toml)
- Interactive setup and management using [Cobra](https://github.com/spf13/cobra) and [Survey](https://github.com/AlecAivazis/survey)
- Commands to set up, update, and delete tasks

## Usage


### Task actions

kasher task create
kasher task update
kasher task delete
kasher task clearAll
kasher task list
kasher task debug




## Configuration

Configuration is stored in:

```
~/Library/Application Support/kasher/config.toml
```

Example config:

```toml
[tasks.qdt]
command = "echo 'Hello, world!'"
cache_duration = "10m"
```

## Dev

* Ensure [Go](https://go.dev/dl/) is installed (version 1.24 or later)
* `$ git clone https://github.com/iammatthew2/kasher`
* `$ cd kasher`
* `$ go run . <optionsHere>`
* `$ go run . task create`
* `$ go run . task list`
new