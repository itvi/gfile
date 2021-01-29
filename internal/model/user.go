package model

import (
	"database/sql"
	e "gfile/pkg/error"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	SN             string
	Name           string
	Email          string
	Password       string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Create(u *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return err
	}

	q := `INSERT INTO app_user(sn,name,email,hashed_password)
		VALUES(?,?,?,?)`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.SN, u.Name, u.Email, string(hashedPassword))
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return e.ErrDuplicate
		}
	}
	return err
}

// GetUsers get all users from database.
func (m *UserModel) GetUsers() ([]*User, error) {
	rows, err := m.DB.Query("SELECT id,sn,name,email FROM app_user;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*User{}
	for rows.Next() {
		u := &User{}
		if err := rows.Scan(&u.ID, &u.SN, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (m *UserModel) Delete(id int) error {
	stmt, err := m.DB.Prepare("DELETE FROM app_user WHERE id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (m *UserModel) Edit(u *User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.HashedPassword), 12)
	if err != nil {
		return err
	}
	q := `UPDATE app_user SET sn=?,name=?,email=?,hashed_password=?
        WHERE id=?`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.SN, u.Name, u.Email, string(hashedPassword), u.ID)
	return err
}

// GetUser method fetch details for a specific user
func (m *UserModel) GetUser(id int) (*User, error) {
	q := "SELECT id,sn,name,email FROM app_user WHERE id=?"
	u := &User{}

	err := m.DB.QueryRow(q, id).Scan(&u.ID, &u.SN, &u.Name, &u.Email)
	if err == sql.ErrNoRows {
		return nil, e.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return u, nil
}

// Authenticate verify where a user exist with the user sn and password
// This will return the relevant user struct
func (m *UserModel) Authenticate(sn, password string) (*User, error) {
	user := &User{}
	row := m.DB.QueryRow(`SELECT id,sn,hashed_password FROM app_user WHERE sn=?`, sn)
	err := row.Scan(&user.ID, &user.SN, &user.HashedPassword)
	if err == sql.ErrNoRows {
		return nil, e.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, e.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}

	return user, nil
}
