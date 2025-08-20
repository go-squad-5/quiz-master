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
	"fmt"
	"strings"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

//	var quizRepo = &quizRepository{
//	  db: sqlmock.New(),
//	}
//
// func NewMySQLRepository(dsn string) (Repository, error) {
func Test_repo_GetAllQuestionByTopic_successfulCase(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	topic := "topic"

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
		WithArgs(topic).
		WillReturnRows(rows)

	repo := &quizRepository{
		db: db,
	}

	_, err = repo.GetAllQuestionByTopic("topic")
	if err != nil {
		t.Errorf("failed: %v", err)
	}
}

func Test_repo_GetAllQuestionByTopic_unableToQuery(t *testing.T) {
	// getquestionbyid
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	topic := "topic"
	execptedError := "Failed to execute query"

	mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
		WithArgs(topic).WillReturnError(fmt.Errorf("%s", execptedError))

	repo := &quizRepository{
		db: db,
	}

	_, err = repo.GetAllQuestionByTopic("topic")
	if err == nil {
		t.Errorf("wanted error: %s, got no err", execptedError)
	} else if err.Error() != execptedError {
		t.Errorf("execpted error: %s, got: %s.", execptedError, err.Error())
	}
}

func Test_repo_GetAllQuestionByTopic_invalidDBResponse(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	topic := "topic"

	rows := sqlmock.NewRows([]string{"id", "question", "options"}).
		AddRow(nil, 14, `["a","b","c","d"]`)

	mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
		WithArgs(topic).
		WillReturnRows(rows)

	repo := &quizRepository{
		db: db,
	}

	_, err = repo.GetAllQuestionByTopic("topic")
	if err == nil {
		t.Errorf("expected sql error, got no err")
	} else if !strings.HasPrefix(err.Error(), "sql") {
		t.Errorf("execpted sql error, got: %s.", err.Error())
	}
}

func Test_repo_GetAllQuestionByTopic_invalidDBJson(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	topic := "topic"
	expectedError := "failed to unmarshal options:"

	rows := sqlmock.NewRows([]string{"id", "question", "options"}).
		AddRow("1", "ques", "a")

	mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
		WithArgs(topic).
		WillReturnRows(rows)

	repo := &quizRepository{
		db: db,
	}

	_, err = repo.GetAllQuestionByTopic("topic")
	if err == nil {
		t.Errorf("expected unmarshal error, got no err")
	} else if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("execpted error format: %s err, got: %s.", expectedError, err.Error())
	}
}

func Test_repo_GetAllQuestionByTopic_rowsIterationErrorCheck(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	topic := "topic"
	expectedError := "error from rows.Err()"

	rows := sqlmock.NewRows([]string{"id", "question", "options"}).
		AddRow("1", "ques", "a")

	mock.ExpectQuery("SELECT id, question, options FROM questions WHERE topic = ?").
		WithArgs(topic).
		WillReturnRows(rows).
		WillReturnError(fmt.Errorf("%s", expectedError))

	repo := &quizRepository{
		db: db,
	}

	_, err = repo.GetAllQuestionByTopic("topic")
	if err == nil {
		t.Errorf("expected error: %s, got no err", expectedError)
	} else if !strings.HasPrefix(err.Error(), expectedError) {
		t.Errorf("execpted error: %s, got: %s.", expectedError, err.Error())
	}
}

// func (r *quizRepository) GetAllQuestionByTopic(topic string) ([]models.Question, error)
// func (r *quizRepository) GetQuestionsByIds(ids []string) ([]models.Question, error)
// func (r *quizRepository) CreateQuiz(ssid string, questions []models.Question) error
// func (r *quizRepository) StoreAnswers(ssid, questionId, answer string, is_correct bool) error
