package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"github.com/shynggys9219/greenlight/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

type UserInfoModel struct {
	DB *sql.DB
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintext string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintext
	p.hash = hash
	return nil
}

func (p *password) Matches(plaintext string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintext))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (m UserInfoModel) Insert(info *UserInfo) error {
	query := "INSERT INTO user_info( name, surname, email, password_hash, role, activated) VALUES ($1,$2,$3, $4, $5, $6) RETURNING id, created_at, version"

	log.Println("inserted to db")

	args := []any{info.Name, info.Surname, info.Email, info.PasswordHash, info.Role, info.Activated}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&info.ID, &info.CreatedAt, &info.Version)
	if err != nil {
		switch {
		case err.Error() == "pq: duplicate key value violates unique constraint 'users_email_key'":
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (m UserInfoModel) GetByEmail(email string) (*UserInfo, error) {
	query := `SELECT id, created_at, name, surname, email, password_hash, role, activated, version FROM user_info WHERE email = $1`

	var info UserInfo

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&info.ID,
		&info.CreatedAt,
		&info.Name,
		&info.Surname,
		&info.Email,
		&info.PasswordHash,
		&info.Role,
		&info.Activated,
		&info.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &info, nil
}

func (m UserInfoModel) GetByID(id int64) (*UserInfo, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, name, surname, email, password_hash, role, activated, version FROM user_info WHERE id = $1`

	var info UserInfo

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&info.ID,
		&info.CreatedAt,
		&info.Name,
		&info.Surname,
		&info.Email,
		&info.PasswordHash,
		&info.Role,
		&info.Activated,
		&info.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &info, nil
}

func (m UserInfoModel) GetAll() []*UserInfo {
	query := "SELECT * FROM user_info ORDER BY id"

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil
	}

	var infos []*UserInfo
	for rows.Next() {
		info := &UserInfo{}
		err = rows.Scan(
			&info.ID,
			&info.Name,
			&info.Surname,
			&info.Email,
			&info.Role,
			&info.Activated,
			&info.CreatedAt,
			&info.UpdatedAt,
			&info.Version,
		)
		if err != nil {
			return nil
		}
		infos = append(infos, info)
	}

	if err = rows.Err(); err != nil {
		return nil
	}
	return infos
}

func (m UserInfoModel) Update(info *UserInfo) error {
	query := "UPDATE user_info SET updated_at = now(), name = $1, surname = $2, email = $3, role = $4, activated = $5, version = version + 1 WHERE id = $6 RETURNING version"

	args := []any{
		info.Name,
		info.Surname,
		info.Email,
		info.Role,
		info.Activated,
		info.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&info.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

func (m UserInfoModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := "DELETE FROM user_info WHERE id = $1"

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateUser(v *validator.Validator, user *UserInfo) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)
	if user.PasswordHash.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.PasswordHash.plaintext)
	}
	if user.PasswordHash.hash == nil {
		panic("missing password hash for user")
	}
}

func (m UserInfoModel) GetForToken(tokenScope, tokenPlaintext string) (*UserInfo, error) {
	// Calculate the SHA-256 hash of the plaintext token provided by the client.
	// Remember that this returns a byte *array* with length 32, not a slice.
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))
	// Set up the SQL query.
	query := `
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
FROM users
INNER JOIN tokens
ON users.id = tokens.user_id
WHERE tokens.hash = $1
AND tokens.scope = $2
AND tokens.expiry > $3`
	// Create a slice containing the query arguments. Notice how we use the [:] operator
	// to get a slice containing the token hash, rather than passing in the array (which
	// is not supported by the pq driver), and that we pass the current time as the
	// value to check against the token expiry.
	args := []any{tokenHash[:], tokenScope, time.Now()}
	var user UserInfo
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Execute the query, scanning the return values into a User struct. If no matching
	// record is found we return an ErrRecordNotFound error.
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	// Return the matching user.
	return &user, nil
}
