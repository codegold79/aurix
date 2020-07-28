package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v4"
)

type dbClient struct {
	*pgx.Conn
}

type ClickCount struct {
	Date  string
	Count int
}

type ClicksView struct {
	ClickCounts []ClickCount
}

func main() {
	ctx := context.Background()

	dbc, err := newDb(ctx)
	if err != nil {
		fmt.Printf("connecting to db: %v", err)
		os.Exit(1)
	}
	defer dbc.Close(ctx)

	if err := dbc.createTable(ctx); err != nil {
		fmt.Printf("creating table: %v", err)
	}

	tmpl := template.Must(
		template.ParseFiles("frontend-template.html"),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		y, m, d := now.Date()
		today := fmt.Sprintf("%v %v, %v", m, d, y)

		// Button has been pushed.
		if r.Method == http.MethodPost {
			dbc.upsertClicksToday(ctx, today, 1)
		}

		// Display what's in click table in database.
		cv, err := dbc.clicksView(ctx)
		if err != nil {
			fmt.Printf("retrieving dates and click counts: %v", err)
		}

		err = tmpl.ExecuteTemplate(w, "frontend-template.html", cv)
		if err != nil {
			fmt.Printf("executing frontend template: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Serving...")
	http.ListenAndServe(":8080", nil)
}
