# gotodocli

`gotodocli` is a web server that wraps calls to [`gotodo`](https://github.com/saifulwebid/gotodo) package.

## Dependency

Make sure that [`gotodo`](https://github.com/saifulwebid/gotodo) dependencies are already fulfilled.

## Installation

Clone this repository:

```sh
mkdir -p $GOPATH/src/github.com/saifulwebid
cd $GOPATH/src/github.com/saifulwebid
git clone https://github.com/saifulwebid/gotodoapp.git
```

Then, build:

```sh
cd $GOPATH/src/github.com/saifulwebid/gotodoapp
go build ./app/gotodocli
```

Executable will be `gotodocli`.

Setup environment variables by either exporting variables to current shell or creating `.env` file. Template to do this is `env.sample`.

Make sure database is already set up.

## Usage

### `./gotodocli getall`

This command returns all Todos stored in database.

Optionally, you can append a `done` argument with either `true` or `false` value (i.e. `./gotodocli getall --done=true`), to get either finished or pending Todos.

### `./gotodocli get [id]`

This command returns a Todo with specified `id`.

### `./gotodocli create`

This command creates a Todo supplied in this arguments:

* `--title`: title of Todo
* `--description`: description of Todo

### `./gotodocli edit [id]`

This command modifies a Todo with values supplied in request body. This command uses `./gotodocli create` arguments.

Only attributes supplied in the arguments will be modified.

### `./gotodocli done [id]`

This command marks a Todo as done.

### `./gotodocli delete [id]`

This command deletes a Todo with a specific `id`.

### `./gotodocli delete-finished`

This command deletes all finished Todos.
