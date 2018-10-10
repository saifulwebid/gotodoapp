# gotodoserver

`gotodoserver` is a web server that wraps calls to [`gotodo`](https://github.com/saifulwebid/gotodo) package.

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
go build ./app/gotodoserver
```

Executable will be `gotodoserver`.

Setup environment variables by either exporting variables to current shell or creating `.env` file. Template to do this is `env.sample`.

Make sure database is already set up.

Then: `./gotodoserver`. It will listen on `http://localhost:PORT`, where `PORT` value is configured on your environment variable.

## Usage

### GET `/`

This endpoint returns all Todos stored in database.

Optionally, you can append a `done` query string with either `true` or `false` value (i.e. `/?done=true`), to get either finished or pending Todos.

### GET `/:id`

This endpoint returns a Todo with specified `id`.

### POST `/`

This endpoint creates a Todo supplied in request body.

Format of a Todo accepted in this endpoint:

```json
{
    "title": "...",
    "description": "..."
}
```

### PATCH `/:id`

This endpoint modifies a Todo with values supplied in request body.

Because this is a `PATCH` endpoint, only attributes supplied in request body will be modified.

This endpoint uses POST `/` request body format.

### PUT `/:id/done`

This endpoint marks a Todo as done. This endpoint accepts no request body.

### DELETE `/:id`

This endpoint deletes a Todo with a specific `id`.

### DELETE `/?done=true`

This endpoint deletes all finished Todos.
