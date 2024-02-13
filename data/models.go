package data

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

var db *sql.DB

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		Transactions: Transactions{},
		Client:       Client{},
		Statement:    Statement{},
		Balance:      Balance{},
	}
}

// Models is the type for this package. Note that any model that is included as a member
// in this type is available to us throughout the application, anywhere that the
// app variable is used, provided that the model is also added in the New function.
type Models struct {
	Transactions Transactions
	Client       Client
	Statement    Statement
	Balance      Balance
}

// User is the structure which holds one user from the database.
type Transactions struct {
	Value       int       `json:"valor"`
	Type        string    `json:"tipo"`
	Description string    `json:"descricao,omitempty"`
	Done_at     time.Time `json:"done_at"`
	ID          int       `json:"client_id"`
}

type Client struct {
	ID      int `json:"id"`
	Limit   int `json:"limit"`
	Balance int `json:"balance"`
}

type Statement struct {
	Balance_details   Balance        `json:"statement_balance"`
	Last_transactions []Transactions `json:"last_transactions"`
}

type Balance struct {
	Total          int       `json:"total"`
	Statement_date time.Time `json:"statement_date"`
	Limit          int       `json:"limite"`
}

func (app Models) GetExtractHandler(clientId int) (*Statement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, _ := db.QueryContext(ctx, "SELECT saldo, limite, now() FROM clients WHERE id = $1", clientId)
	var balance Balance
	err := rows.Scan(
		balance.Limit,
		balance.Statement_date,
		balance.Total,
	)

	if err != nil {
		return nil, err
	}

	rows, _ = db.QueryContext(ctx, "SELECT valor, tipo, descricao, done_at, client_id FROM transactions WHERE client_id = $1 ORDER BY id DESC LIMIT 10", clientId)

	fmt.Println(rows)
	defer rows.Close()

	var transactions []Transactions

	for rows.Next() {
		var transaction Transactions
		err := rows.Scan(
			&transaction.Description,
			&transaction.Done_at,
			&transaction.ID,
			&transaction.Type,
			&transaction.Value,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		transactions = append(transactions, transaction)

		fmt.Println(transactions)
	}

	result := Statement{
		Balance_details:   balance,
		Last_transactions: transactions,
	}

	fmt.Println(result)

	return &result, nil
}
