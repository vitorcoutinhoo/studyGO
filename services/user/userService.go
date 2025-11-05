package user

import (
	"acommerce-api/types/models"
	"database/sql"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	rows, err := s.db.Query("SELECT * FROM user EHERE email = ?", email)

	if err != nil {
		return nil, err
	}

	u := new(models.User)

	for rows.Next() {
		u, err = scanRowInUser(rows)

		if err != nil {
			return nil, err
		}
	}

	if u.ID == 0 {
		return nil, fmt.Errorf("User not found")
	}

	return u, nil
}

func scanRowInUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
