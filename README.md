# Golang Store API

Production-ready E-commerce backend, REST API in [Go](https://go.dev/) using [JWT authentication](https://jwt.io/introduction) &amp; [MySQL](https://www.w3schools.com/MySQL/default.asp).

## API documentation

**Note:** for auth guarded endpoints, you have to hit the `/login` endpoint and retrieve the JWT token. Then, you can use that token string and place it as a query param `?token=` or in the `Authorization` header.

## Auth

| Method | Endpoint    | Description                       | Request Body                           | Response                      | Authentication |
| ------ | ----------- | --------------------------------- | -------------------------------------- | ----------------------------- | -------------- |
| POST   | `/login`    | Logs in a user and returns a JWT. | Email and password                     | 200 OK / 400 Bad Request      | No             |
| POST   | `/register` | Registers a new user.             | First name, last name, email, password | 201 Created / 400 Bad Request | No             |

## Users

| GET | `/users` | Retrieves a list of all users. | N/A | 200 OK / 500 Internal Server Error | Yes (Admin) |
| GET | `/users/{id}` | Retrieves a user by their ID. | User ID | 200 OK / 400 Bad Request / 500 Internal Server Error | Yes (Admin) |

## Error Responses

- **400 Bad Request**: This response is returned if the request body is invalid, or if required parameters are missing.
- **401 Unauthorized**: This response is returned if the JWT token is missing or invalid (for protected routes).
- **500 Internal Server Error**: This response is returned if there is an issue with the server processing the request.

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
