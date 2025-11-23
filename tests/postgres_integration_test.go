package tests

import (
	"context"
	"testing"
	"time"

	"github.com/makson2134/go-qa-service/internal/repository/postgres"
	"github.com/pressly/goose/v3"
	testcontainerspostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(t *testing.T) (*postgres.DB, func()) {
	ctx := context.Background()

	postgresContainer, err := testcontainerspostgres.Run(ctx,
		"postgres:17-alpine",
		testcontainerspostgres.WithDatabase("testdb"),
		testcontainerspostgres.WithUsername("testuser"),
		testcontainerspostgres.WithPassword("testpass"),
		testcontainerspostgres.BasicWaitStrategies(),
		testcontainerspostgres.WithSQLDriver("pgx"),
	)
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}

	connStr, err := postgresContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	db, err := postgres.New(connStr, 10, 5, time.Hour)
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.GetDB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("failed to set goose dialect: %v", err)
	}

	if err := goose.Up(sqlDB, "../migrations"); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	cleanup := func() {
		if err := db.Close(); err != nil {
			t.Logf("failed to close database connection: %v", err)
		}
		if err := postgresContainer.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup
}

func TestMultipleAnswersFromSameUser(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	question, err := db.Create("What is Go?")
	if err != nil {
		t.Fatalf("failed to create question: %v", err)
	}

	userID := "user123"

	_, err = db.CreateAnswer(question.ID, userID, "Go is a programming language")
	if err != nil {
		t.Fatalf("failed to create first answer: %v", err)
	}

	_, err = db.CreateAnswer(question.ID, userID, "Go was created by Google")
	if err != nil {
		t.Fatalf("failed to create second answer: %v", err)
	}

	fetchedQuestion, err := db.GetByID(question.ID)
	if err != nil {
		t.Fatalf("failed to get question with answers: %v", err)
	}

	if len(fetchedQuestion.Answers) != 2 {
		t.Fatalf("expected 2 answers, got %d", len(fetchedQuestion.Answers))
	}

	for _, answer := range fetchedQuestion.Answers {
		if answer.UserID != userID {
			t.Errorf("expected user_id %s, got %s", userID, answer.UserID)
		}
	}
}

func TestCascadeDeleteQuestionDeletesAnswers(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	question, err := db.Create("What is Docker?")
	if err != nil {
		t.Fatalf("failed to create question: %v", err)
	}

	answer1, err := db.CreateAnswer(question.ID, "user1", "Docker is a containerization platform")
	if err != nil {
		t.Fatalf("failed to create first answer: %v", err)
	}

	answer2, err := db.CreateAnswer(question.ID, "user2", "Docker uses containers")
	if err != nil {
		t.Fatalf("failed to create second answer: %v", err)
	}

	if err := db.Delete(question.ID); err != nil {
		t.Fatalf("failed to delete question: %v", err)
	}

	_, err = db.GetAnswerByID(answer1.ID)
	if err == nil {
		t.Error("expected answer1 to be deleted, but it still exists")
	}

	_, err = db.GetAnswerByID(answer2.ID)
	if err == nil {
		t.Error("expected answer2 to be deleted, but it still exists")
	}
}
