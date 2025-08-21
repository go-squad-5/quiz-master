# Quiz Master

The Quiz master manages the quiz creation and scoring quiz.
## Prerequisties
- A running mysql server, the default name is quizdb, You can check all default config in `internal/config/config.go`. (you can name database anything you want and give dsn in env, it is also recommended)
- Seed the mysql database using the `scripts/db.sql` (`db.sql` contains table creation and seeding commands)

## Setup
- To run the quiz manager, execute the following commands
```bash
cd path/to/quiz-master
go mod tidy
go run cmd/web/main.go
```
-   you can to set up the environment variables, example declared in `.env.example`:
```bash
cp .env.example .env
```
Then, edit the `.env` file to set your environment variables.

- To build the quiz master, you can use the following command:
```bash
go build -o <your-build-file-path> ./cmd/web/main.go
```
Then you can run the quiz client: `<your-build-file-path>`

## Config
### You can setup config information in `internal/config/config.go`
- change worker count
- change port

## Run Tests

- To run the tests for the quiz master, you can use the following command:

```bash
go test -cover ./...
```

- For better test coverage, you can run the tests with the following command:

```bash
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

- Using `tparse` you can parse the test results and generate a report:

**Note** You need to install tparse to use it. To install tparse see the docs [here](https://github.com/mfridman/tparse). Also make sure `~/go/bin` is add to environment paths.
```bash
set -o pipefail && go test -cover -coverprofile=coverage.out -json ./... | tparse -all
```

### For detailed report you can checkout [here](https://go-squad-5.github.io/quiz-master/#file0)
