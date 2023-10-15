# Mastermind REST API
REST API that simulates the role of Mastermind's codemaker.

[Mastermind](https://en.wikipedia.org/wiki/Mastermind_(board_game)) is a code-breaking game for two players. This API simulates the role of codemaker. As a codebreaker, you can guess the code by sending a four digit number to the codemaker where each digit is between 1 to 6. The API will respond with a list. This list will be empty in case non of the digits were guessed correctly or filled with a combination of ones and zeros for correctly guessed digits. One indicates that a digit has the correct position, and zero that it doesn't.

## HTTP tests

```bash
GIN_MODE=release go test -v
```
