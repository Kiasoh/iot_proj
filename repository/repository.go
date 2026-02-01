package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"iot_proj/models"
)

type Repository struct {
	DB *pgxpool.Pool
}

func (r *Repository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (email, password, key_card, access_level)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(
		context.Background(),
		query,
		user.Email,
		user.Password,
		user.KeyCard,
		user.AccessLevel,
	).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *Repository) GetUser(keyCard string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, key_card, access_level, last_entered, created_at
		FROM users
		WHERE key_card = $1
	`
	err := r.DB.QueryRow(context.Background(), query, keyCard).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.KeyCard,
		&user.AccessLevel,
		&user.LastAccessed,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, key_card, access_level, last_entered, created_at
		FROM users
		WHERE email = $1
	`
	err := r.DB.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.KeyCard,
		&user.AccessLevel,
		&user.LastAccessed,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) GetUserByID(id int) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password, key_card, access_level, last_entered, created_at
		FROM users
		WHERE id = $1
	`
	err := r.DB.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.KeyCard,
		&user.AccessLevel,
		&user.LastAccessed,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *Repository) CreateEntryLog(entryLog *models.EntryLog) error {
	query := `
		INSERT INTO entry_logs (key_card, status, message)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`
	err := r.DB.QueryRow(
		context.Background(),
		query,
		entryLog.KeyCard,
		entryLog.Status,
		entryLog.Message,
	).Scan(&entryLog.ID, &entryLog.CreatedAt)
	return err
}

func (r *Repository) GetEntryLogs(keyCard string) ([]*models.EntryLog, error) {
	rows, err := r.DB.Query(
		context.Background(),
		`
			SELECT id, key_card, status, message, created_at
			FROM entry_logs
			WHERE key_card = $1
			ORDER BY created_at DESC
		`,
		keyCard,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.EntryLog
	for rows.Next() {
		log := &models.EntryLog{}
		var message *string
		err := rows.Scan(
			&log.ID,
			&log.KeyCard,
			&log.Status,
			&message,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		log.Message = message
		logs = append(logs, log)
	}

	return logs, nil
}
func (r *Repository) GetAllEntryLogsPaginated(limit, offset int) ([]*models.EntryLog, error) {
	rows, err := r.DB.Query(
		context.Background(),
		`
			SELECT id, key_card, status, message, created_at
			FROM entry_logs
			ORDER BY created_at DESC
			LIMIT $1 OFFSET $2
		`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.EntryLog
	for rows.Next() {
		log := &models.EntryLog{}
		var message *string
		err := rows.Scan(
			&log.ID,
			&log.KeyCard,
			&log.Status,
			&message,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		log.Message = message
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *Repository) GetEntryLogsPaginated(keyCard string, limit, offset int) ([]*models.EntryLog, error) {
	rows, err := r.DB.Query(
		context.Background(),
		`
			SELECT id, key_card, status, message, created_at
			FROM entry_logs
			WHERE key_card = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3
		`,
		keyCard,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*models.EntryLog
	for rows.Next() {
		log := &models.EntryLog{}
		var message *string
		err := rows.Scan(
			&log.ID,
			&log.KeyCard,
			&log.Status,
			&message,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		log.Message = message
		logs = append(logs, log)
	}

	return logs, nil
}

func (r *Repository) UpdateUserLastAccessed(keyCard string) error {
	query := `
		UPDATE users
		SET last_entered = $1
		WHERE key_card = $2
	`
	_, err := r.DB.Exec(context.Background(), query, time.Now(), keyCard)
	return err
}

func (r *Repository) GetUsers(limit, offset int) ([]*models.User, error) {
	rows, err := r.DB.Query(
		context.Background(),
		`
			SELECT id, email, password, key_card, access_level, last_entered, created_at
			FROM users
			ORDER BY id
			LIMIT $1 OFFSET $2
		`,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.KeyCard,
			&user.AccessLevel,
			&user.LastAccessed,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) UpdateUserKeyCard(userID int, keyCard string) error {
	query := `
		UPDATE users
		SET key_card = $1
		WHERE id = $2
	`
	_, err := r.DB.Exec(context.Background(), query, keyCard, userID)
	return err
}

func (r *Repository) UpdateUserAccessLevel(userID int, accessLevel int) error {
	query := `
		UPDATE users
		SET access_level = $1
		WHERE id = $2
	`
	_, err := r.DB.Exec(context.Background(), query, accessLevel, userID)
	return err
}
