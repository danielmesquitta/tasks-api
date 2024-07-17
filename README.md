# Tasks API

This is an Tasks API CRUD.
It is written in Go and uses the Echo framework and a MySQL database.

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
