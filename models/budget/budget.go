package budget

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
	Year     int
	Month    Month
	incomes  []float64
	expenses []Expense
}

type Expense struct {
	category string
	amount   []float64
}

func CreateBudget(year int, month Month) Budget {

	expenses := []Expense{}

	return Budget{
		Year:     year,
		Month:    month,
		incomes:  []float64{},
		expenses: expenses,
	}
}
