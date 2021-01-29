package model

import (
	"database/sql"
	e "gfile/pkg/error"
	"strings"
)

type Role struct {
	ID          int
	Name        string
	Description string
}

type RoleModel struct {
	DB *sql.DB
}

func (m *RoleModel) Create(r *Role) error {
	q := `INSERT INTO app_role(name,description) VALUES(?,?);`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.Name, r.Description)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return e.ErrDuplicate
		}
	}
	return err
}

func (m *RoleModel) GetRoles() ([]*Role, error) {
	q := `SELECT id,name,description FROM app_role;`
	rows, err := m.DB.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	app_roles := []*Role{}
	for rows.Next() {
		r := &Role{}
		if err := rows.Scan(&r.ID, &r.Name, &r.Description); err != nil {
			return nil, err
		}
		app_roles = append(app_roles, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return app_roles, nil
}

func (m *RoleModel) GetRole(id int) (*Role, error) {
	q := `SELECT id,name,description FROM app_role
          WHERE id=?;`
	r := &Role{}

	err := m.DB.QueryRow(q, id).Scan(&r.ID, &r.Name, &r.Description)
	if err == sql.ErrNoRows {
		return nil, e.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return r, nil
}

func (m *RoleModel) Delete(id int) error {
	q := `DELETE FROM app_role WHERE id=?;`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	return err
}

func (m *RoleModel) Edit(r *Role) error {
	q := `UPDATE app_role SET name=?,description=? WHERE id=?;`
	stmt, err := m.DB.Prepare(q)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(r.Name, r.Description, r.ID)
	return err
}
