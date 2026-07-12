package budget

import (
	"budget-cli/db"
	"fmt"
)

type Month int

const (
	January Month = iota + 1
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func (m Month) String() string {
	names := [...]string{
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	}
	if m < January || m > December {
		return "Unknown"
	}
	return names[m]
}

type Budget struct {
	Year  int
	Month Month
	Items []BudgetItem
}

type BudgetItem struct {
	Category Category
	Name     string
	Amount   float64
}

type Category struct {
	name string
	kind string
}

func CreateBudget(year int, month Month) Budget {

	return Budget{
		Year:  year,
		Month: month,
		Items: nil,
	}
}

// buildBudgetItem constructs a BudgetItem from item data and category data
func buildBudgetItem(itemName string, itemAmount float64, categoryName string, categoryKind string) BudgetItem {
	return BudgetItem{
		Name:   itemName,
		Amount: itemAmount,
		Category: Category{
			name: categoryName,
			kind: categoryKind,
		},
	}
}

func GetBudget(month Month, year int) (*Budget, error) {

	query := `SELECT bi.name, bi.planned_amount_cents,
		c.name, c.kind
		FROM budgets b
		INNER JOIN budget_items bi ON b.id = bi.budget_id
		INNER JOIN categories c ON bi.category_id = c.id
		WHERE b.month = ? AND b.year = ?
		`

	rows, err := db.DB.Query(query, month, year)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	budget := Budget{
		Items: []BudgetItem{},
		Month: month,
		Year:  year,
	}

	for rows.Next() {
		var itemName string
		var itemAmount float64
		var categoryName string
		var categoryKind string

		err := rows.Scan(&itemName, &itemAmount, &categoryName, &categoryKind)
		if err != nil {
			return nil, err
		}
		budget.Items = append(budget.Items, buildBudgetItem(itemName, itemAmount, categoryName, categoryKind))
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &budget, nil
}

func GetAllBudgets() ([]*Budget, error) {
	query := `SELECT b.year, b.month, bi.name, bi.planned_amount_cents,
		c.name, c.kind
		FROM budgets b
		LEFT JOIN budget_items bi ON b.id = bi.budget_id
		LEFT JOIN categories c ON bi.category_id = c.id
		ORDER BY b.year, b.month
		`

	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Use a map to group items by budget key, and track insertion order
	budgetMap := make(map[string]*Budget)
	budgetKeys := []string{}

	for rows.Next() {
		var year int
		var month int
		var itemName *string
		var itemAmount *float64
		var categoryName *string
		var categoryKind *string

		err := rows.Scan(&year, &month, &itemName, &itemAmount, &categoryName, &categoryKind)
		if err != nil {
			return nil, err
		}

		budgetKey := fmt.Sprintf("%d-%d", year, month)

		// Create budget if it doesn't exist yet
		if _, exists := budgetMap[budgetKey]; !exists {
			budgetMap[budgetKey] = &Budget{
				Year:  year,
				Month: Month(month),
				Items: []BudgetItem{},
			}
			budgetKeys = append(budgetKeys, budgetKey)
		}

		// Add item to budget if it exists (LEFT JOIN may produce null items)
		if itemName != nil && itemAmount != nil && categoryName != nil && categoryKind != nil {
			budgetMap[budgetKey].Items = append(
				budgetMap[budgetKey].Items,
				buildBudgetItem(*itemName, *itemAmount, *categoryName, *categoryKind),
			)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Reconstruct slice in the order budgets were encountered (which is sorted by year, month)
	result := make([]*Budget, 0, len(budgetKeys))
	for _, key := range budgetKeys {
		result = append(result, budgetMap[key])
	}

	return result, nil
}
