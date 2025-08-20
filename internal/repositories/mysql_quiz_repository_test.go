package repositories

// unable to query
// -- get questions by topic
// -- git question by id

// unable to exec
// -- create quiz
// -- stroe answer

// invalid db response
// -- create quiz
// -- store answer
// -- get questions by topic
// -- git question by id

// invalid db json
// -- get questions by topic
// -- git question by id
import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-squad-5/quiz-master/internal/models"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func Test_repo_GetAllQuestionByTopic(t *testing.T) {
	args := []driver.Value{"topic"}
	tests := []struct {
		name        string
		prepareMock func(sqlmock.Sqlmock, error)
		wantErr     bool
		errorPrefix string
	}{
		{
			name: "get all question by topic: successful case",
			prepareMock: func(mock sqlmock.Sqlmock, _ error) {
				rows := sqlmock.NewRows([]string{"id", "question", "options"}).
					AddRow("1", "ques", `["a","b","c","d"]`).
					AddRow("2", "ques", `["a","b","c","d"]`).
					AddRow("3", "ques", `["a","b","c","d"]`).
					AddRow("4", "ques", `["a","b","c","d"]`).
					AddRow("5", "ques", `["a","b","c","d"]`).
					AddRow("6", "ques", `["a","b","c","d"]`).
					AddRow("7", "ques", `["a","b","c","d"]`).
					AddRow("8", "ques", `["a","b","c","d"]`).
					AddRow("9", "ques", `["a","b","c","d"]`)

				mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
					WithArgs(args...).
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "get question by topic: unable to query",
			prepareMock: func(mock sqlmock.Sqlmock, expectedError error) {
				mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
					WithArgs(args...).WillReturnError(fmt.Errorf("%s", expectedError))
			},
			wantErr:     true,
			errorPrefix: "Failed to execute query",
		},
		{
			name: "get question by topic: invalid db response",
			prepareMock: func(mock sqlmock.Sqlmock, _ error) {
				rows := sqlmock.NewRows([]string{"id", "question", "options"}).
					AddRow(nil, 14, `["a","b","c","d"]`)

				mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
					WithArgs(args...).
					WillReturnRows(rows)
			},
			wantErr:     true,
			errorPrefix: "sql",
		},
		{
			name: "get question by topic: invalid json",
			prepareMock: func(mock sqlmock.Sqlmock, _ error) {
				rows := sqlmock.NewRows([]string{"id", "question", "options"}).
					AddRow("1", "ques", "a")

				mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
					WithArgs(args...).
					WillReturnRows(rows)
			},
			wantErr:     true,
			errorPrefix: "failed to unmarshal options:",
		},
	}

	for _, test := range tests {
		func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.prepareMock(mock, errors.New(test.errorPrefix))

			repo := &quizRepository{
				db: db,
			}

			_, err = repo.GetAllQuestionByTopic("topic")
			if test.wantErr {
				if err == nil {
					t.Errorf("test: %s\n wanted error, got no error\n\n", test.name)
				} else if !strings.HasPrefix(err.Error(), test.errorPrefix) {
					t.Errorf("test: %s\n wanted error of format:\t %s\n got err:\t\t %s\n\n", test.name, test.errorPrefix, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("test: %s\n wanted no errors, got: %s\n\n", test.name, err.Error())
				}
			}
		}()
	}
}

func Test_repo_GetQuestionsByIds(t *testing.T) {
	tests := []struct {
		name        string
		prepareMock func(sqlmock.Sqlmock, []string, error)
		args        []string
		wantErr     bool
		errorPrefix string
	}{
		{
			name: "successful case",
			prepareMock: func(mock sqlmock.Sqlmock, args []string, _ error) {
				driverValueArgs := make([]driver.Value, 0)
				for _, v := range args {
					driverValueArgs = append(driverValueArgs, driver.Value(v))
				}
				rows := sqlmock.NewRows([]string{"id", "answer"}).
					AddRow("1", "a").
					AddRow("1", "a").
					AddRow("1", "a").
					AddRow("1", "a").
					AddRow("1", "a")

				mock.ExpectQuery("SELECT id, answer FROM questions WHERE id IN ").WithArgs(driverValueArgs...).WillReturnRows(rows)
			},
			args:    []string{"id1", "id2", "id3"},
			wantErr: false,
		},
		{
			name: "no ids in args",
			prepareMock: func(mock sqlmock.Sqlmock, args []string, err error) {
				driverValueArgs := make([]driver.Value, 0)
				for _, v := range args {
					driverValueArgs = append(driverValueArgs, driver.Value(v))
				}
				mock.ExpectQuery("SELECT id, answer FROM questions WHERE id IN ").WithArgs(driverValueArgs...)
			},
			args:        []string{},
			wantErr:     true,
			errorPrefix: "no question IDs provided",
		},
		{
			name: "unable to query",
			prepareMock: func(mock sqlmock.Sqlmock, args []string, _ error) {
				driverValueArgs := make([]driver.Value, 0)
				for _, v := range args {
					driverValueArgs = append(driverValueArgs, driver.Value(v))
				}
				mock.ExpectQuery("SELECT id, answer FROM questions WHERE id IN ").WithArgs(driverValueArgs...).WillReturnError(errors.New("query failed for mocking"))
			},
			args:        []string{"id1", "id2", "id3"},
			wantErr:     true,
			errorPrefix: "query failed: ",
		},
		{
			name: "invalid db response",
			prepareMock: func(mock sqlmock.Sqlmock, args []string, _ error) {
				driverValueArgs := make([]driver.Value, 0)
				for _, v := range args {
					driverValueArgs = append(driverValueArgs, driver.Value(v))
				}
				rows := sqlmock.NewRows([]string{"id", "answer"}).
					AddRow(nil, nil)

				mock.ExpectQuery("SELECT id, answer FROM questions WHERE id IN ").WithArgs(driverValueArgs...).WillReturnRows(rows)
			},
			args:        []string{"id1", "id2", "id3"},
			wantErr:     true,
			errorPrefix: "failed to scan question: ",
		},
	}

	for _, test := range tests {
		func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.prepareMock(mock, test.args, errors.New(test.errorPrefix))

			repo := &quizRepository{
				db: db,
			}

			_, err = repo.GetQuestionsByIds(test.args)
			if test.wantErr {
				if err == nil {
					t.Errorf("test: %s\n wanted error, got no error\n\n", test.name)
				} else if !strings.HasPrefix(err.Error(), test.errorPrefix) {
					t.Errorf("test: %s\n wanted error of format:\t %s\n got err:\t\t %s\n\n", test.name, test.errorPrefix, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("test: %s\n wanted no errors, got: %s\n\n", test.name, err.Error())
				}
			}
		}()
	}
}

func Test_repo_CreateQuiz(t *testing.T) {
	// args := {"ssid", "ques_id"}

	tests := []struct {
		name        string
		prepareMock func(sqlmock.Sqlmock, error)
		wantErr     bool
		err         string
	}{
		{
			name: "successful case",
			prepareMock: func(mock sqlmock.Sqlmock, _ error) {
				mock.ExpectExec("INSERT INTO quizzes").WithArgs("ssid", "ques_id").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
			err:     "",
		},
		{
			name: "unable to exec query",
			prepareMock: func(mock sqlmock.Sqlmock, err error) {
				mock.ExpectExec("INSERT INTO quizzes").WithArgs("ssid", "ques_id").WillReturnError(err)
			},
			wantErr: true,
			err:     "mock error",
		},
	}

	for _, test := range tests {
		func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.prepareMock(mock, errors.New(test.err))

			repo := &quizRepository{
				db: db,
			}

			err = repo.CreateQuiz("ssid", []models.Question{
				{
					Id: "ques_id",
				},
			})
			if test.wantErr {
				if err == nil {
					t.Errorf("test: %s\n wanted error, got no error\n\n", test.name)
				} else if !strings.HasPrefix(err.Error(), test.err) {
					t.Errorf("test: %s\n wanted error of format:\t %s\n got err:\t\t %s\n\n", test.name, test.err, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("test: %s\n wanted no errors, got: %s\n\n", test.name, err.Error())
				}
			}
		}()
	}
}

func Test_repo_StoreAnswer(t *testing.T) {
	tests := []struct {
		name        string
		prepareMock func(sqlmock.Sqlmock)
		wantErr     bool
		err         string
	}{
		{
			name: "successful response",
			prepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE quizzes SET answer").WithArgs("answer", false, "ssid", "ques_id").WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "unable to exec",
			prepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec("UPDATE quizzes SET answer").WithArgs("answer", false, "ssid", "ques_id").WillReturnError(errors.New("mock error"))
			},
			wantErr: true,
			err:     "failed to store answer: mock error",
		},
	}

	for _, test := range tests {
		func() {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			test.prepareMock(mock)

			repo := &quizRepository{
				db: db,
			}

			err = repo.StoreAnswers("ssid", "ques_id", "answer", false)
			if test.wantErr {
				if err == nil {
					t.Errorf("test: %s\n wanted error, got no error\n\n", test.name)
				} else if !strings.HasPrefix(err.Error(), test.err) {
					t.Errorf("test: %s\n wanted error of format:\t %s\n got err:\t\t %s\n\n", test.name, test.err, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("test: %s\n wanted no errors, got: %s\n\n", test.name, err.Error())
				}
			}
		}()
	}
}
