package db

import (
	"context"
	"database/sql"
	"fmt"
)

const tableName = "test"

const selectSQL = "SELECT name, age FROM %s WHERE age > $1"

type Repository struct {
	selectStmt *sql.Stmt
	db         *sql.DB
}

type People struct {
	Name string
	Age  int
}

func NewRepository(db *sql.DB) (*Repository, error) {
	selectStmt, err := db.Prepare(fmt.Sprintf(selectSQL, tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to prepare stmt: %w", err)
	}

	return &Repository{
		selectStmt: selectStmt,
		db:         db,
	}, nil
}

func (r *Repository) Select(ctx context.Context, age int) ([]People, error) {
	rows, err := r.selectStmt.QueryContext(ctx, age)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	people := make([]People, 0)

	for rows.Next() {
		var p People
		err := rows.Scan(
			&p.Name,
			&p.Age,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan a row: %w", err)
		}

		people = append(people, p)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed while iterating over rows: %w", rows.Err())
	}

	return people, nil
}
