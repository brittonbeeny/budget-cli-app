package budget

import (
	"database/sql"
	"testing"
	"time"

	"budget-cli/db"

	_ "github.com/mattn/go-sqlite3"
)

// setupTestDB creates and initializes an in-memory SQLite database with the same schema as the real database
func setupTestDB(t *testing.T) *sql.DB {
	testDB, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory database: %v", err)
	}

	schema := `
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

	if _, err := testDB.Exec(schema); err != nil {
		t.Fatalf("failed to execute schema: %v", err)
	}

	return testDB
}

// seedTestData populates the test database with sample data
func seedTestData(t *testing.T, testDB *sql.DB) {
	now := time.Now()

	// Insert categories
	categories := []struct {
		name string
		kind string
	}{
		{"Groceries", "expense"},
		{"Utilities", "expense"},
		{"Salary", "income"},
		{"Rent", "expense"},
	}

	categoryIDs := make(map[string]int)
	for _, cat := range categories {
		var id int
		err := testDB.QueryRow(
			"INSERT INTO categories (name, kind, created_at) VALUES (?, ?, ?) RETURNING id",
			cat.name, cat.kind, now,
		).Scan(&id)
		if err != nil {
			t.Fatalf("failed to insert category: %v", err)
		}
		categoryIDs[cat.name] = id
	}

	// Insert budget
	var budgetID int
	err := testDB.QueryRow(
		"INSERT INTO budgets (year, month, created_at) VALUES (?, ?, ?) RETURNING id",
		2024, 5, now,
	).Scan(&budgetID)
	if err != nil {
		t.Fatalf("failed to insert budget: %v", err)
	}

	// Insert budget items
	budgetItems := []struct {
		category string
		amount   int
		name     string
	}{
		{"Groceries", 50000, "Weekly groceries"},
		{"Utilities", 15000, "Electric and water"},
		{"Rent", 120000, "Monthly rent"},
	}

	for _, item := range budgetItems {
		_, err := testDB.Exec(
			"INSERT INTO budget_items (budget_id, category_id, planned_amount_cents, name, created_at) VALUES (?, ?, ?, ?, ?)",
			budgetID, categoryIDs[item.category], item.amount, item.name, now,
		)
		if err != nil {
			t.Fatalf("failed to insert budget item: %v", err)
		}
	}
}

// seedMultipleBudgets populates the test database with multiple budgets for testing GetAllBudgets
func seedMultipleBudgets(t *testing.T, testDB *sql.DB) {
	now := time.Now()

	// Insert categories
	categories := []struct {
		name string
		kind string
	}{
		{"Groceries", "expense"},
		{"Utilities", "expense"},
		{"Salary", "income"},
		{"Rent", "expense"},
	}

	categoryIDs := make(map[string]int)
	for _, cat := range categories {
		var id int
		err := testDB.QueryRow(
			"INSERT INTO categories (name, kind, created_at) VALUES (?, ?, ?) RETURNING id",
			cat.name, cat.kind, now,
		).Scan(&id)
		if err != nil {
			t.Fatalf("failed to insert category: %v", err)
		}
		categoryIDs[cat.name] = id
	}

	// Insert multiple budgets with items
	budgets := []struct {
		year  int
		month int
	}{
		{2024, 3},
		{2024, 5},
		{2024, 12},
	}

	for _, b := range budgets {
		var budgetID int
		err := testDB.QueryRow(
			"INSERT INTO budgets (year, month, created_at) VALUES (?, ?, ?) RETURNING id",
			b.year, b.month, now,
		).Scan(&budgetID)
		if err != nil {
			t.Fatalf("failed to insert budget: %v", err)
		}

		// Add items to this budget
		budgetItems := []struct {
			category string
			amount   int
			name     string
		}{
			{"Groceries", 50000, "Weekly groceries"},
			{"Utilities", 15000, "Electric and water"},
		}

		for _, item := range budgetItems {
			_, err := testDB.Exec(
				"INSERT INTO budget_items (budget_id, category_id, planned_amount_cents, name, created_at) VALUES (?, ?, ?, ?, ?)",
				budgetID, categoryIDs[item.category], item.amount, item.name, now,
			)
			if err != nil {
				t.Fatalf("failed to insert budget item: %v", err)
			}
		}
	}
}

// TestGetBudget tests the GetBudget function with in-memory database
func TestGetBudget(t *testing.T) {
	// Setup test database
	testDB := setupTestDB(t)
	defer testDB.Close()

	// Seed test data
	seedTestData(t, testDB)

	// Replace the global DB with our test DB
	oldDB := db.DB
	db.DB = testDB
	defer func() { db.DB = oldDB }()

	// Test GetBudget for existing budget
	budget, err := GetBudget(May, 2024)
	if err != nil {
		t.Fatalf("GetBudget failed: %v", err)
	}

	if budget == nil {
		t.Fatal("expected budget to not be nil")
	}

	if budget.Year != 2024 {
		t.Errorf("expected year 2024, got %d", budget.Year)
	}

	if budget.Month != May {
		t.Errorf("expected month May, got %v", budget.Month)
	}

	if len(budget.Items) != 3 {
		t.Errorf("expected 3 budget items, got %d", len(budget.Items))
	}

	// Verify budget items
	expectedItems := map[string]struct {
		amount   float64
		category string
	}{
		"Weekly groceries":   {amount: 50000, category: "Groceries"},
		"Electric and water": {amount: 15000, category: "Utilities"},
		"Monthly rent":       {amount: 120000, category: "Rent"},
	}

	for _, item := range budget.Items {
		expected, exists := expectedItems[item.Name]
		if !exists {
			t.Errorf("unexpected budget item: %s", item.Name)
			continue
		}

		if item.Amount != expected.amount {
			t.Errorf("item %s: expected amount %f, got %f", item.Name, expected.amount, item.Amount)
		}

		if item.Category.name != expected.category {
			t.Errorf("item %s: expected category %s, got %s", item.Name, expected.category, item.Category.name)
		}
	}

	// Test GetBudget for non-existing budget
	emptyBudget, err := GetBudget(January, 2024)
	if err != nil {
		t.Fatalf("GetBudget for non-existent budget failed: %v", err)
	}

	if emptyBudget == nil {
		t.Fatal("expected budget to not be nil even when empty")
	}

	if len(emptyBudget.Items) != 0 {
		t.Errorf("expected 0 budget items for non-existent budget, got %d", len(emptyBudget.Items))
	}
}

// TestGetAllBudgets tests the GetAllBudgets function with in-memory database
func TestGetAllBudgets(t *testing.T) {
	// Setup test database
	testDB := setupTestDB(t)
	defer testDB.Close()

	// Seed test data with multiple budgets
	seedMultipleBudgets(t, testDB)

	// Replace the global DB with our test DB
	oldDB := db.DB
	db.DB = testDB
	defer func() { db.DB = oldDB }()

	// Test GetAllBudgets
	allBudgets, err := GetAllBudgets()
	if err != nil {
		t.Fatalf("GetAllBudgets failed: %v", err)
	}

	if allBudgets == nil {
		t.Fatal("expected budgets slice to not be nil")
	}

	if len(allBudgets) != 3 {
		t.Errorf("expected 3 budgets, got %d", len(allBudgets))
	}

	// Verify budgets are sorted by year and month
	expectedBudgets := []struct {
		year  int
		month Month
		items int
	}{
		{2024, March, 2},
		{2024, May, 2},
		{2024, December, 2},
	}

	for i, expected := range expectedBudgets {
		if i >= len(allBudgets) {
			t.Errorf("not enough budgets returned")
			break
		}

		budget := allBudgets[i]
		if budget.Year != expected.year {
			t.Errorf("budget %d: expected year %d, got %d", i, expected.year, budget.Year)
		}

		if budget.Month != expected.month {
			t.Errorf("budget %d: expected month %v, got %v", i, expected.month, budget.Month)
		}

		if len(budget.Items) != expected.items {
			t.Errorf("budget %d: expected %d items, got %d", i, expected.items, len(budget.Items))
		}
	}
}
