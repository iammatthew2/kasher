# kasher

**kasher** is a CLI tool for defining, running, and caching shell commands. This util is intended to reduce the wait time for simple fetch commands that are slow, rate-limited, or rarely change (getting Kubernetes pod info (`kubectl get pods`), for example)

Kasher stores cachable shell commands as "tasks". All task actions (create, update, delete, ...) are available as subcommands of `task` (`kasher task list`, `kasher task help`, `kasher task create`)

Once a task is stored it can be run directly on `kasher`, like so `kasher <taskNameHere>`

Tasks are assigned an expiration (specify how many hours, minutes or seconds until cache is stale). When a task is executed before the expiration is hit, the cache is used. When a task is executed after then a new request is maade.

## Instalation

- tbd

## Features

- Define named tasks that wrap shell commands
- Cache task output for a set expiration time
- Interactive task definition flow via a set of commands to set up and modify tasks - see full list under **Task actions**

## Basic usage

You have a common request you run that takes a little too long or you want access to the response when you're offline (`aws s3 ls`, for example). The data is fairly static.

Capture the shell command as a task:

`$ kasher task create`

provide data at the prompts:

- Name: myTask
- Shell command: aws s3 ls 
- Cache expiration: 48h (I could choose a larger number, but I want sort of fresh data)

Now call this task as often as you need:

`$ kasher myTask`

Kasher will reduce the number of requests you make but will keep the data available to you.

## Advanced usage

Force the cache to refresh with `--force` (`-f`), no matter the expiration:

`$ kasher myTask -f`

Pass a shell command directly to kasher using `createFor`:

`$ kasher createFor cat README.md`



## Available task actions

- `kasher task create [name]` — interactively create a new task (optionally specify the name as an argument)
- `kasher task createFor <command>` — create a new task for a given shell command (provide the command as arguments)
- `kasher task update` — update an existing task
- `kasher task delete` — delete a task
- `kasher task clearAll` — delete all tasks/settings
- `kasher task list` — list all tasks
- `kasher task debug` — debug output


## Dev

Interested in fixing/improving Kasher?

* Ensure [Go](https://go.dev/dl/) is installed (version 1.24 or later)
* `$ git clone https://github.com/iammatthew2/kasher`
* `$ cd kasher`
* `$ go run . <optionsHere>`
* `$ go run . task create`
* `$ go run . task list`
