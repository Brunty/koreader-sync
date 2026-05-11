package dao

import (
	"database/sql"
	"errors"

	"github.com/brunty/koreader-sync-server/db"
	"github.com/brunty/koreader-sync-server/types"
)

func SelectProgress(userId int64, document string) (*types.Progress, error) {
	var progress = types.Progress{}
	progress.UserID = userId
	progress.Document = document

	query := "SELECT id, progress, percentage, device, device_id, timestamp FROM progress WHERE user_id = ? AND document = ?"
	err := db.DBCon.QueryRow(query, userId, document).Scan(
		&progress.ID,
		&progress.Progress,
		&progress.Percentage,
		&progress.Device,
		&progress.DeviceID,
		&progress.Timestamp,
	)

	if err != nil {
		// if there's no rows, that's fine, just return nil nil
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &progress, nil
}

func StoreProgress(progress types.Progress) error {
	args := []interface{}{
		sql.Named("user_id", progress.UserID),
		sql.Named("document", progress.Document),
		sql.Named("progress", progress.Progress),
		sql.Named("percentage", progress.Percentage),
		sql.Named("device", progress.Device),
		sql.Named("device_id", progress.DeviceID),
		sql.Named("timestamp", progress.Timestamp),
	}

	_, err := db.DBCon.Exec(`
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
		return err
	}

	return nil
}
