package repositories

import (
	"context"
	"database/sql"

	"github.com/sgatu/ezmail/internal/domain/models/auth"
	"github.com/uptrace/bun"
)

type mysqlSessionRepository struct {
	db *bun.DB
}

func NewMysqlSessionRepository(connection *bun.DB) *mysqlSessionRepository {
	return &mysqlSessionRepository{
		db: connection,
	}
}

func (repo *mysqlSessionRepository) GetSessionById(ctx context.Context, id string) (*auth.Session, error) {
	session := &auth.Session{Id: id}
	err := repo.db.NewSelect().Model(session).WherePK().Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, auth.ErrNoSessionFound
	} else if err != nil {
		return nil, err
	}
	return session, nil
}

func (repo *mysqlSessionRepository) GetSessionBySessionId(ctx context.Context, sessionId string) (*auth.Session, error) {
	session := &auth.Session{}
	err := repo.db.NewSelect().Model(session).Where("session_id = ?", sessionId).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, auth.ErrNoSessionFound
	} else if err != nil {
		return nil, err
	}
	return session, nil
}

func (repo *mysqlSessionRepository) Save(ctx context.Context, session *auth.Session) error {
	_, err := repo.db.NewInsert().Model(session).Exec(ctx)
	return err
}
