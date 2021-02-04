package jsonstruct

type ForexExchange struct {
	Rates struct {
		CHF float64 `json:"CHF"`
		USD float64 `json:"USD"`
	} `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}
