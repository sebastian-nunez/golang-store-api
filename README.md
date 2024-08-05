# Golang Store API

Production-ready E-commerce backend, REST API in [Go](https://go.dev/) using [JWT authentication](https://jwt.io/introduction) &amp; [MySQL](https://www.w3schools.com/MySQL/default.asp).

## Features

TODO(sebastian-nunez): add feature list

## Getting started

### Running locally

Make sure to have Go 1.22+ and Make installed, and then run:

```bash
make run
```

_The project requires environment variables to be set. You can find the list of required variables in the `.env.template` file._

### Database migrations

We are using [golang-migrate](https://github.com/golang-migrate/migrate/tree/master) to ease all database migrations.

#### Create a migration

To create a new database migration, run:

```bash
make migration <migration-name>
```

Then, you can find the create (up) and teardown (down) scripts in `/cmd/migrate/migrations`.

#### Applying all database migration

To apply all existing database migrations, run:

```bash
make migrate-up
```

#### Turning down database migrations

To remove all database migrations, run:

```bash
make migrate-down
```
