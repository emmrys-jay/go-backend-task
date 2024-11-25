package main

type AccountStatement struct {
	Name           string `json:"name,omitempty"`
	Currency       string `json:"currency,omitempty"`
	CurrencySymbol string `json:"currency_symbol,omitempty"`
	Iban           []Iban `json:"iban,omitempty"`
	BalanceSummary struct {
		Products []Product `json:"products,omitempty"`
	} `json:"balance_summary,omitempty"`
	Transactions []Transaction `json:"transactions,omitempty"`
	CreatedAt    string        `json:"created_at,omitempty"`
	Address      Address       `json:"address,omitempty"`
}

type Address struct {
	HouseNo int    `json:"house_no,omitempty"`
	Street  string `json:"street,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	Country string `json:"country,omitempty"`
}

type Iban struct {
	No  string `json:"no,omitempty"`
	Bic string `json:"bic,omitempty"`
}

type Product struct {
	Product        string  `json:"product,omitempty"`
	OpeningBalance float64 `json:"opening_balance,omitempty"`
	MoneyIn        float64 `json:"money_in,omitempty"`
	MoneyOut       float64 `json:"money_out,omitempty"`
	ClosingBalance float64 `json:"closing_balance,omitempty"`
}

type Transaction struct {
	Date        string  `json:"date,omitempty"`
	Description string  `json:"description,omitempty"`
	MoneyIn     float64 `json:"money_in,omitempty"`
	MoneyOut    float64 `json:"money_out,omitempty"`
	Balance     float64 `json:"balance,omitempty"`
}
