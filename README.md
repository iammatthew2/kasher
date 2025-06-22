# kasher

**Kasher** is a CLI tool for defining, running, and caching shell commands. This util is intended to reduce the wait time for simple fetch commands that are slow, rate-limited, or rarely change (getting Kubernetes pod info (`kubectl get pods`), for example)

Kasher stores cachable shell commands as "tasks". All task actions (create, update, list, ...) are available as subcommands of `task` (`kasher task create`, `kasher task update`, `kasher task lists`)

Once a task is stored it can be run directly on `kasher`: `kasher <taskNameHere>`

Tasks are assigned an expiration (specify how many hours, minutes or seconds until cache is stale). When a task is executed before the expiration is hit, the cache is used. When a task is executed after expiration is hit then a new request is maade.

## Instalation

- tbd

## Features

- Define named tasks that wrap shell commands
- Cache task output for a set expiration time (see [ParseDuration](https://pkg.go.dev/time#ParseDuration)
- Interactive task definition flow via a set of commands to set up and modify tasks - see full list under **Task actions**

## Usage

### Basic

Lets say you have a common request you run that takes a little too long or you want access to the response when you're offline (`aws s3 ls`, for example).

Capture the shell command as a task:

`$ kasher task create`

provide data at the prompts:

    Name: myTask
    Shell command: aws s3 ls 
    Expiration: 48h (see help for more time options)

Now call this task as often as you need:

`$ kasher myTask`

Kasher will reduce the number of requests you make but will keep the data available to you.

### Fuzzy search for tasks

Run `kasher` without any args to trigger the fuzzu search task finder: `$ kasher`


### More options

Force the cache to refresh with `--force` (`-f`), no matter the expiration: `$ kasher myTask -f`

Pass a shell command directly to kasher with *createFor: `$ kasher createFor cat README.md`

View all saved tasks: `$ kasher task list`

## Available task actions

- `kasher task create [name]` — interactively create a new task (optionally specify the name as an argument)
- `kasher task createFor <command>` — create a new task for a given shell command (provide the command as arguments)
- `kasher task update` — update an existing task
- `kasher task delete` — delete a task
- `kasher task clearAll` — delete all tasks/settings
- `kasher task list` — list all tasks


## Dev

Interested in fixing/improving Kasher?

* Ensure [Go](https://go.dev/dl/) is installed (version 1.24 or later)
* `$ git clone https://github.com/iammatthew2/kasher`
* `$ cd kasher`
* `$ go run . <optionsHere>`
* `$ go run . task create`
* `$ go run . task list`
