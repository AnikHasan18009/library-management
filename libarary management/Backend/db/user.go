package db

import (
	"database/sql"
	"fmt"
	"library-service/logger"
	"log"
	"log/slog"
	"net/http"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

type LoggedUser struct {
	Email    string `json:"email" db:"email" validate:"required,max=50,email"`
	Password string `json:"password" db:"password" validate:"required,min=8,max=50"`
}

type RegisteredUser struct {
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
type NewUser struct {
	RegisteredUser
	Role     string `json:"role" db:"role"`
	Approved bool   `json:"approved" db:"approved"`
}
type ApprovedUser struct {
	Name  string `json:"name"  db:"name"`
	Email string `json:"email" db:"email"`
}
type SignedUser struct {
	Id    int    `json:"id"  db:"id"`
	Email string `json:"email" db:"email"`
}

type UserRepo struct {
	UserRepoTable string
}

var userRepo *UserRepo

func InitUserRepo() {
	userRepo = &UserRepo{
		UserRepoTable: `"USER"`,
	}

}

func GenerateHashedPassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func GetUserRepo() *UserRepo {

	return userRepo
}
func (ur *UserRepo) VerifyUserCredentials(user LoggedUser) (int, error) {
	var (
		err          error
		passwordInDB string
		id           int
	)
	Db := GetReadDB()
	queryBulder := GetQueryBuilder().Select("PASSWORD", "ID").From("USER").Where(sq.Eq{"EMAIL": user.Email})

	query, args, err := queryBulder.ToSql()

	if err != nil {
		slog.Error("error during select query creation", logger.Extra(map[string]any{
			"error": err.Error(),
		}))
		return 0, err
	}

	err = Db.QueryRow(query, args...).Scan(&passwordInDB, &id)
	if err != nil {
		slog.Error("error querying database", logger.Extra(map[string]any{
			"error": err.Error(),
		}))
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordInDB), []byte(user.Password))
	if err != nil {
		return 0, fmt.Errorf("wrong password")
	}

	return id, nil

}

func (ur *UserRepo) getInsertionQueryAndArgsForUser(newUser NewUser) (string, []any, error) {
	queryBuilder := GetQueryBuilder().
		Insert(ur.UserRepoTable).
		Columns("name", "email", "password", "role", "approved").
		Values(newUser.Name, newUser.Email, newUser.Password, newUser.Role, newUser.Approved).
		Suffix(`
	on conflict (email) do update set
	password = excluded.password,
	name = excluded.name
	returning id
	`)
	return queryBuilder.ToSql()

}

func (ur *UserRepo) Insert(newUser NewUser) (int, error) {
	var (
		approved bool
		id       int
		err      error
	)
	Db := GetWriteDB()

	if approved, err = ur.CheckIfUserIsActive(newUser.Email); approved {
		err = fmt.Errorf("email already registered")
		slog.Error("email already registerd", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newUser,
		}))
		return http.StatusBadRequest, err
	}

	if err != nil && err != sql.ErrNoRows {
		return http.StatusInternalServerError, err
	}

	if newUser.Password, err = GenerateHashedPassword(newUser.Password); err != nil {
		slog.Error("Error hashing password", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newUser,
		}))
		return http.StatusInternalServerError, err
	}

	query, args, err := ur.getInsertionQueryAndArgsForUser(newUser)
	fmt.Println(query)

	if err != nil {
		slog.Error("Error creating  query", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newUser,
		}))
		return http.StatusInternalServerError, err
	}

	err = Db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		slog.Error("Failed registration", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": newUser,
		}))

		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil

}

// Approve Reader
func (ur *UserRepo) getReaderApprovalQuery(email string) (string, []any, error) {
	queryBuilder := GetQueryBuilder().
		Update(ur.UserRepoTable).
		Set("approved", true).Where(sq.Eq{"email": email}).
		Suffix(`returning approved`)

	return queryBuilder.ToSql()

}

func (ur *UserRepo) ApproveReader(email string) (int, error) {
	var (
		err      error
		approved bool
	)
	Db := GetWriteDB()

	query, args, err := ur.getReaderApprovalQuery(email)
	fmt.Println(query)

	if err != nil {
		slog.Error("Error creating approval query", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": email,
		}))
		return http.StatusInternalServerError, err
	}

	err = Db.QueryRow(query, args...).Scan(&approved)
	if err != nil {
		slog.Error("Failed approval", logger.Extra(map[string]any{
			"error":   err.Error(),
			"payload": email,
		}))
		if err == sql.ErrNoRows {

			return http.StatusBadRequest, err
		}
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil

}

func (ur *UserRepo) CheckIfUserIsActive(email string) (bool, error) {
	var (
		approved bool
		err      error
	)
	Db := GetReadDB()
	qb := GetQueryBuilder()

	queryBuilder := qb.Select("approved").From(ur.UserRepoTable).Where(sq.Eq{"email": email})
	query, args, err := queryBuilder.ToSql()
	fmt.Println(query)

	if err != nil {
		slog.Error(
			"Failed to create query for user active check",
			logger.Extra(map[string]any{
				"error": err.Error()}))
		return approved, err
	}

	if err = Db.QueryRow(query, args...).Scan(&approved); err != nil && err != sql.ErrNoRows {

		slog.Error("error during user existance check", logger.Extra(map[string]any{
			"error": err.Error()}))
	}

	return approved, nil
}

func DeleteUserByIdFromDB(id int) error {
	Db := GetWriteDB()
	_, err := Db.Exec("delete from users where id = $1", id)
	return err

}

func CountOfRecordsBasedOnConditions(conditions map[string]string) int {
	var count int
	Db := GetReadDB()
	query := "select count(name) from users where "

	for k, v := range conditions {
		switch {
		case k == "id":
			query += k + " = " + v + " and "
		default:
			query += k + " = '" + v + "' and "
		}

	}
	query += "true"
	fmt.Println(query)
	err := Db.QueryRow(query).Scan(&count)

	if err != nil {
		log.Fatal(err)
	}
	return count
}
