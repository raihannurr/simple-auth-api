package repository_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/errors"
	"github.com/raihannurr/simple-auth-api/internal/repository"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
)

func TestMysqlRepository(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer mockDB.Close()

	mock.ExpectQuery("SELECT VERSION()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("8.0.35"))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs("user1", 1).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "verified", "description"}).AddRow(1, "user1", "user1@example.com", false, ""))
	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "verified", "description"}).AddRow(1, "user1", "user1@example.com", false, ""))

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs(2, 1).WillReturnError(errors.New("user not found"))
	mock.ExpectQuery("SELECT \\* FROM `users`").WithArgs("user2", 1).WillReturnError(errors.New("user not found"))

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `users`").WillReturnError(errors.New("user not found"))
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WillReturnError(errors.New("unique constraint violation"))
	mock.ExpectRollback()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `users`").WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()

	repo := repository.NewMysqlRepository(config.DatabaseConfig{
		Connection: mysql.New(mysql.Config{
			Conn: mockDB,
		}),
		GormConfig: gorm.Config{
			NowFunc: func() time.Time {
				return time.Time{}
			},
		},
	})

	user, err := repo.CreateUser("user1", "user1@example.com", "password")
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, uint(1), user.ID)

	user, err = repo.GetUserByUsername("user1")
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, uint(1), user.ID)

	user, err = repo.GetUserByID(1)
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, uint(1), user.ID)

	assert.Equal(t, "", user.Description)
	user, err = repo.UpdateUser(1, entity.User{ID: 1, Description: "user1 description"})
	assert.Nil(t, err)
	assert.Equal(t, "user1 description", user.Description)

	_, err = repo.GetUserByID(2)
	assert.Equal(t, errors.ErrUserNotFound, err)

	_, err = repo.GetUserByUsername("user2")
	assert.Equal(t, errors.ErrUserNotFound, err)

	_, err = repo.UpdateUser(2, entity.User{ID: 2, Description: "user2 description"})
	assert.Equal(t, errors.ErrUserNotFound, err)

	_, err = repo.CreateUser("user1", "", "")
	assert.NotNil(t, err)

	userNew, err := repo.CreateUser("user-new", "user-new@example.com", "password")
	assert.Nil(t, err)
	assert.Equal(t, "user-new", userNew.Username)
	assert.Equal(t, uint(2), userNew.ID)
}
