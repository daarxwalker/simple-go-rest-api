# Simple Go REST API

## Prerequisites
- Unix-based OS
- [Taskfile](https://taskfile.dev/)
- [Docker (docker compose)](https://www.docker.com/)
- [Tilt](https://tilt.dev/)
- [golang-migrate](https://formulae.brew.sh/formula/golang-migrate)

<br>

## How to run the project?
- ``cd bin/development``
- ``tilt up``

<br>

## How to stop the project?
- ``CTRL + C``
- ``tilt down``

<br>

## How to test project?
``task test``

<br>

## How to create a migration?
``task migrate:create NAME=<nazev>``

<br>

## How to migrate?
``task migrate:up ENV_PATH=bin/development/.env``

## Endpoints
- ``GET http://localhost/api/v1/users/:id``
- ``POST http://localhost/api/v1/users``