# kasher

**kasher** is a Go-based CLI tool for defining, running, and caching named shell tasks. It helps you wrap shell commands as reusable tasks, cache their output for a configurable expiration time, and manage them interactively.

## Features

- Define named tasks that wrap shell commands
- Cache task output for a set expiration time
- Store configuration (location based on OS and user setup, MacOS: ~/Library/Application Support/kasher/config.toml)
- Interactive setup and management using [Cobra](https://github.com/spf13/cobra) and [Survey](https://github.com/AlecAivazis/survey)
- Commands to set up, update, and delete tasks

## Installation

Clone the repository and build:

```sh
git clone https://github.com/yourusername/kasher.git
cd kasher
go build -o kasher
```

## Usage

### Add a Task

```sh
kasher setup
```
Follow the interactive prompts to define a new task.

### Update or Delete a Task

```sh
kasher update
kasher delete
```
(Commands may be implemented as the project evolves.)

### Run a Task

(Planned feature: run a named task and use cached output if valid.)

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

## Roadmap

...