package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

func newDb(ctx context.Context) (dbClient, error) {
	connURL := fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_HOST"),
		"postgres",
	)

	conn, err := pgx.Connect(ctx, connURL)
	if err != nil {
		return dbClient{}, err
	}

	return dbClient{conn}, nil
}

func (c dbClient) createTable(ctx context.Context) error {
	_, err := c.Exec(ctx, "CREATE TABLE IF NOT EXISTS clicks(date text PRIMARY KEY, count integer)")
	if err != nil {
		return err
	}
	return nil
}

func (c dbClient) clicksView(ctx context.Context) (ClicksView, error) {
	var cv ClicksView
	var cc ClickCount

	rows, err := c.Query(ctx, "SELECT date, count FROM clicks ORDER BY date DESC LIMIT 25")
	if err != nil {
		return cv, err
	}

	for rows.Next() {
		err := rows.Scan(&cc.Date, &cc.Count)
		if err != nil {
			return cv, err
		}

		cv.ClickCounts = append(cv.ClickCounts, cc)
	}

	return cv, nil
}

func (c dbClient) upsertClicksToday(ctx context.Context, today string, count int) error {
	qry := `
		INSERT INTO clicks (date, count) 
		VALUES ($1, $2)
		ON CONFLICT (date) DO UPDATE 
		SET count=clicks.count+1
	`

	_, err := c.Exec(ctx, qry, today, count)
	if err != nil {
		return err
	}

	return err
}
