# Tasks API

## Overview

This is an Tasks API CRUD.
It is written in Go and uses the Echo framework and a MySQL database.
It implements unit tests, dependency injection, clean architecture concepts
and uses the repository, builder and option patterns.

## Dependencies

1. [Go 1.22.4](https://go.dev/doc/install)
2. [Make](https://www.gnu.org/software/make/)
3. [Docker](https://www.docker.com/) with a MySQL container or [MySQL](https://www.mysql.com/downloads/) itself

## Running the API

1. Clone the repository.
2. Install the dependencies with `make install`.
3. Copy `.env.example` to `.env` and fill in your environment variables.
4. Run the migrations with `make migrations_up`
5. Finally, start the application with live-reload enabled running `make dev`.

## Tests and coverage

Unit tests can be run with `make test`, and coverage with `make coverage`.

## Swagger documentation

All API endpoints are documented in the Swagger available at the `/api/v1/docs/index.html` route.

## Key notes

- User passwords are hashed with the bcrypt algorithm, so passwords cannot be decrypted
- Task summary are hashed with the AES algorithm, so they can be decrypted (it is hashed due to security, since it is known that the summary can contain personal information)
- There are user roles to define resource permissions
- There is validation in the input data in every use case
