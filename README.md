# taskManager

[![codecov](https://codecov.io/gh/brionac626/taskManager/graph/badge.svg?token=vKI2foy3CH)](https://codecov.io/gh/brionac626/taskManager)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/brionac626/deck-rng/main/LICENSE)

This is a backend server project for manage user's tasks

## How to install

>Make sure you installed Go in your local machines

1. clone the project first

```sh
git clone git@github.com:brionac626/taskManager.git
```

2. build the project by Go command line tool

```sh
go build -o app main.go
```

or build it by `make` command

```sh
make build
```

3. run the binary file from your terminal

```sh
./app server
```

or run it by `make` command

```sh
make run
```

## How to build docker image for the project

Build docker image by docker command line tool

```sh
docker build -f ./Dockerfile -t task-manager:latest .
```

or build it by `make` command

```sh
make docker-build
```

## How to run the application as a docker container

>You need to build the docker image first

Run the application by docker command line tool

```sh
docker run -d -p 8080:8080 --name task-manager task-manager:latest
```

or run it by `make` command

```sh
make docker-run
```

For more details, please refer to the  [makefile](makefile)