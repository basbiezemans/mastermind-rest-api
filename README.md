# Mastermind REST API

REST API that simulates the role of Mastermind's codemaker.

The API generates a secret four-digit code. As a codebreaker, submit guesses
where each digit is between `1` and `6`. Feedback uses:

- `●` for a digit in the correct position
- `○` for a correct digit in the wrong position

## Requirements

- Go 1.25 or newer
- SQLite

## Run the server

```bash
go run main.go
```

The server listens on `http://localhost:8080`. The SQLite database is stored
in `data/mastermind.db`, and application errors are written to
`data/error.log`.

## API

### Get endpoint information

```http
GET /
```

Example:

```bash
curl http://localhost:8080/
```

### Create a game

```http
POST /create
```

Example:

```bash
curl -X POST http://localhost:8080/create
```

Successful responses return `201 Created`:

```json
{
  "message": "A new game has been created. Good luck!",
  "token": "20d245fd-f724-4e1c-a818-04b3dd33ef5d"
}
```

Save the returned token for subsequent requests.

### Get a game

```http
GET /games/:token
```

Example:

```bash
curl http://localhost:8080/games/20d245fd-f724-4e1c-a818-04b3dd33ef5d
```

The secret code is not included in the response.

Possible responses:

- `200 OK` when the game exists
- `400 Bad Request` when the token is not a valid UUID
- `404 Not Found` when the game does not exist

### Submit a guess

```http
PATCH /games/:token
```

Submit the guess as an `application/x-www-form-urlencoded` field named
`guess`. Each guess must contain exactly four digits between `1` and `6`.

Example:

```bash
curl -X PATCH \
  -H "Content-Type: application/x-www-form-urlencoded" \
  --data "guess=1234" \
  http://localhost:8080/games/20d245fd-f724-4e1c-a818-04b3dd33ef5d
```

Possible responses:

- `200 OK` when the guess is processed
- `400 Bad Request` for an invalid token or missing guess
- `404 Not Found` when the game does not exist
- `500 Internal Server Error` for an unexpected server failure

### Delete a game

```http
DELETE /games/:token
```

Example:

```bash
curl -X DELETE \
  http://localhost:8080/games/20d245fd-f724-4e1c-a818-04b3dd33ef5d
```

Successful deletion returns `204 No Content`. An invalid UUID returns
`400 Bad Request`.

## Tests

Run all package tests with:

```bash
go test ./...
```

Run the HTTP-focused test command with verbose output:

```bash
GIN_MODE=test go test -v
```
