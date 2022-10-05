# demo-basic
A basic app with a single device (CRUD with model) "basic" demonstrating how to build
an application with Chassis and Meta.

## Requirements

- [Go](https://go.dev/)
- [Docker](https://www.docker.com/) and [docker-compose](https://docs.docker.com/compose/install/)
- [Meta](https://github.com/oligoden/meta)
- [Graphviz](https://graphviz.org/) (optional)

## Development

Requirements for development are to have Go, Docker and docker-compose installed.
Meta is used to build the service and terraform for deployment.

Clone the repo and copy the dev environment override `meta.override.example.json`
file and rename it to `meta.override.json`. Finally, build the source files with
`meta build`.

```bash
    git clone git@github.com:oligoden/demo-basic.git
    cd demo-basic
    cp meta.override.example.json meta.override.json
    meta build
```

Note: run `meta up` in a separate terminal to update source code while
actively coding.

### Unit testing

Before tests can be run the database must be running

```bash
    docker-compose up -d dev-db
```

To run unit tests run: `go test github.com/oligoden/demo-basic`

### Running locally

The following rules must be added to your `/etc/hosts` file

```
    127.0.0.1 test.com
    127.0.0.1 api.test.com
```

Run with docker-compose.

```bash
    docker-compose up -d dev-db
    docker-compose up basic
```

Go to `test.com:8080`

### Testing the API

Go to `api.test.com:8080` to see the OpenAPI client. Here you can test:

- creating a record (POST /basics)
- viewing records (GET /basics)
- updating a record (PUT /basics/UC)

When testing the PUT /basics/UC end-point remember to replace the default UC value (xy)
in the textbox with a UC value returned by the POST or GET end-points.

Refer to the `oas.json` file for the OAS discription of the API.

### Viewing the development DB

To connect to the DB point your client to:
```
localhost:3310
database: basic
username: user
password: pass
```