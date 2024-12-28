# How to start the project

This is a streaming service project built in Go. It allows streaming multimedia content using AWS for file management and PostgreSQL to handle user and video information. Additionally, it supports Docker to facilitate deployment in local and production environments.

## Tech Stack

**Client:** React, TailwindCSS

**Server:** Golang, Gin, AWS, PostgreSQL

## Requirements

To use this project, you need:

### Mandatory

- **Git**: To clone the repository.
- **Go**: To run and compile the project. (Recommended version: 1.23.1, used in `go.mod`).

### Optional

- **Docker**: To containerize the application.
- **Docker Compose**: To orchestrate services if using Docker.

## Installation

```bash
git clone https://github.com/Unbot2313/go-streaming-service.git
cd go-streaming-service/
```

## Usage/Examples

To use it with Go!:

```bash
go mod tidy
go run main.go
```

With Docker (includes the PostgreSQL instance locally):

```bash
docker compose up --build
```

## Contributing

Contributions are always welcome!

Please star a new fork, then make a pull request.

In case of changing the documentation, do the following:

```bash
go install github.com/go-swagger/go-swagger/cmd/swagger@latest
swag init # update the documentation from the swagger
```

## Features

- Return existing videos and their content
- RefreshTokens
- Handle live streaming
- Add the Auth middleware to the streaming routes (middleware already done)
- Finish the README.md
- Optimize the Dockerfile