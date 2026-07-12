package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

const schema = `
PRAGMA foreign_keys = OFF;

CREATE TABLE IF NOT EXISTS budgets (
id INTEGER PRIMARY KEY,
year INTEGER NOT NULL,
month INTEGER NOT NULL,
created_at DATETIME,
UNIQUE(month,year),
CHECK(month BETWEEN 1 AND 12)
);

CREATE TABLE IF NOT EXISTS categories (
id INTEGER PRIMARY KEY,
name TEXT NOT NULL,
kind TEXT NOT NULL,
parent_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
notes TEXT,
created_at DATETIME NOT NULL,
updated_at DATETIME,
UNIQUE(name, kind),
CHECK(kind IN ('income','expense'))
);

CREATE TABLE IF NOT EXISTS budget_items (
id INTEGER PRIMARY KEY,
budget_id INTEGER NOT NULL REFERENCES budgets(id) ON DELETE CASCADE,
category_id INTEGER NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
planned_amount_cents INTEGER NOT NULL,
name TEXT,
created_at DATETIME NOT NULL,
updated_at DATETIME,
UNIQUE(budget_id, category_id)
);

CREATE INDEX IF NOT EXISTS idx_budgets_year_month ON budgets(year, month);
CREATE INDEX IF NOT EXISTS idx_budget_items_budget ON budget_items(budget_id);
CREATE INDEX IF NOT EXISTS idx_categories_kind ON categories(kind);

PRAGMA foreign_keys = ON;
`

func InitDB() error {
	var err error

	DB, err = sql.Open("sqlite3", "budget.db")

	if err != nil {
		log.Fatalf("failed to open database: %v", err)
		return err
	}

	//Connection pool. Reuse
	DB.SetMaxOpenConns(6)
	//At least N connections open at all times
	DB.SetMaxIdleConns(3)

	if _, err := DB.Exec(schema); err != nil {
		log.Fatal("Could not execute schema")
		return err
	}

	return nil
}
