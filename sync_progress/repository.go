package sync_progress

import (
	"context"
	"database/sql"
	"errors"
)

type SyncProgressRepository interface {
	SelectByUserIDAndDocument(ctx context.Context, userID int64, document string) (*SyncProgress, error)
	Store(ctx context.Context, syncProgress SyncProgress) (*int64, error)
}

type SyncProgressRepositorySQLite struct {
	db *sql.DB
}

func NewSyncProgressRepository(db *sql.DB) SyncProgressRepository {
	return &SyncProgressRepositorySQLite{db: db}
}

func (r *SyncProgressRepositorySQLite) SelectByUserIDAndDocument(ctx context.Context, userID int64, document string) (*SyncProgress, error) {
	var progress = SyncProgress{}
	progress.UserID = userID
	progress.Document = document

	query := "SELECT id, progress, percentage, device, device_id, timestamp FROM progress WHERE user_id = ? AND document = ?"
	err := r.db.QueryRowContext(ctx, query, userID, document).Scan(
		&progress.ID,
		&progress.Progress,
		&progress.Percentage,
		&progress.Device,
		&progress.DeviceID,
		&progress.Timestamp,
	)

	if err != nil {
		// if there are no rows, that's fine, just return nil, nil
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &progress, nil
}

func (r *SyncProgressRepositorySQLite) Store(ctx context.Context, progress SyncProgress) (*int64, error) {
	args := []interface{}{
		sql.Named("user_id", progress.UserID),
		sql.Named("document", progress.Document),
		sql.Named("progress", progress.Progress),
		sql.Named("percentage", progress.Percentage),
		sql.Named("device", progress.Device),
		sql.Named("device_id", progress.DeviceID),
		sql.Named("timestamp", progress.Timestamp),
	}

	res, err := r.db.ExecContext(ctx, `
		INSERT INTO progress (
			user_id,
			document,
			progress,
			percentage,
			device,
			device_id,
			timestamp
		) VALUES (
			@user_id,
			@document,
			@progress,
			@percentage,
			@device,
			@device_id,
			@timestamp
		)
		ON CONFLICT(user_id, document) DO UPDATE SET
			 progress = @progress,
			 percentage = @percentage,
			 device = @device,
			 device_id = @device_id,
			 timestamp = @timestamp
		`,
		args...,
	)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()

	return &id, nil
}
