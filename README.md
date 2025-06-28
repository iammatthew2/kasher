# kasher

**Kasher** is a CLI tool for defining and running shell commands and caching their responses. It’s designed to reduce wait times for simple fetch commands that are slow, rate-limited, or rarely change—such as retrieving Kubernetes pod info (`kubectl get pods`).

Kasher stores cachcacheable able shell commands as "tasks". All task actions (create, update, list, ...) are available as subcommands of *task* (`kasher task create`, `kasher task update`, `kasher task list`, ...)

Once a task is stored it can be run directly on *kasher*: `kasher <taskNameHere>`

Tasks are assigned an expiration (specify how many hours, minutes or seconds until cache is stale). When a task is executed before the expiration is hit, the cache is used. When a task is executed after expiration is hit then a new request is made and cached.

![kasher](https://github.com/user-attachments/assets/e7bcd84e-79a6-46ef-b994-d783a9d20d84)


## Installation

### Brew install

```
brew tap iammatthew2/kasher
brew install kasher
```
> [!NOTE]
> The above tap links to https://github.com/iammatthew2/homebrew-kasher/, which links back here to Kasher

### Alternatives
See latest release in [releases](https://github.com/iammatthew2/kasher/releases) for direct installation instructions


## Features

- Define named tasks that wrap shell commands
- Cache task output for a set expiration time (see [ParseDuration](https://pkg.go.dev/time#ParseDuration))
- Interactive task definition flow via a set of commands to set up and modify tasks - see full list under **Task actions**

## Usage

Lets say you have a common request you frequently run that takes a little too long or you want access to the response when you're offline (`aws s3 ls`, for example).

Capture the shell command as a task:

`$ kasher task create`

Provide data at the prompts:

    Name: myTask
    Shell command: aws s3 ls 
    Cache expiration (e.g. 10m, 1h, 2h30m, 45s, 1.5h, 500ms): 20s
    Notes (optional): just an example

Now call this task as often as you need:

`$ kasher myTask`

Kasher will reduce the number of requests you make while keeping the data available to you.

### Fuzzy search for tasks

Run `kasher` without any args to trigger the fuzzy search task finder: `$ kasher`

### Available task actions

- `kasher task create [name]` — interactively create a new task (optionally specify the name as an argument)
- `kasher task createFor <command>` — create a new task for a given shell command (provide the command as arguments):

  `$ kasher createFor "echo starting && sleep 5 && echo ending"`
- `kasher task update` — update an existing task
- `kasher task delete` — delete a task
- `kasher task clearAll` — delete all tasks/settings
- `kasher task list` — list all tasks

### Flags

`--force` (`-f`), Force the cache to refresh, no matter the expiration: `$ kasher myTask -f`. This will execute the task immediately and cache its response.


## Dev

Interested in fixing/improving Kasher?

* Ensure [Go](https://go.dev/dl/) is installed (version 1.24 or later)
* `$ git clone https://github.com/iammatthew2/kasher`
* `$ cd kasher`
* `$ go run . <optionsHere>`
* `$ go run . task create`
* `$ go run . task list`
