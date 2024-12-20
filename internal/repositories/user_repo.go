package repositories

import (
	"database/sql"
	"errors"

	"github.com/JosephKristian/project-migration/internal/models"
)

type UserRepo struct {
	DB *sql.DB
}

// Save menyimpan user ke database.
func (repo *UserRepo) Save(user *models.User) error {
	query := `INSERT INTO users (name, email, phone, password, uuid, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := repo.DB.Exec(query, user.Name, user.Email, user.Phone, user.Password, user.UUID, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return errors.New("failed to save user: " + err.Error())
	}
	return nil
}

// FindByEmail mencari user berdasarkan email.
func (repo *UserRepo) FindByEmail(email string) (*models.User, error) {
	query := `SELECT name, email, phone, password, uuid, created_at, updated_at 
			  FROM users WHERE email = $1`

	row := repo.DB.QueryRow(query, email)

	var user models.User
	if err := row.Scan(&user.Name, &user.Email, &user.Phone, &user.Password, &user.UUID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Tidak ditemukan
		}
		return nil, errors.New("failed to find user: " + err.Error())
	}

	return &user, nil
}

// CheckEmailExist memeriksa apakah email sudah terdaftar.
func (repo *UserRepo) CheckEmailExist(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	err := repo.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, errors.New("failed to check email existence: " + err.Error())
	}
	return count > 0, nil
}
