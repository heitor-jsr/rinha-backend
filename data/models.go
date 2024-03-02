package data

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const dbTimeout = time.Second * 3

var db *pgxpool.Pool

// New is the function used to create an instance of the data package. It returns the type
// Model, which embeds all the types we want to be available to our application.
func New(dbPool *pgxpool.Pool) Models {
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
	Transactions      Transactions
	Client            Client
	Statement         Statement
	Balance           Balance
	TransactionResult TransactionResult
}

// User is the structure which holds one user from the database.
type Transactions struct {
	ID           int       `json:"-"`
	Value        int       `json:"valor"`
	Type         string    `json:"tipo"`
	Description  string    `json:"descricao"`
	Realizada_em time.Time `json:"realizada_em"`
	Cliente_ID   int       `json:"-"`
}

type Client struct {
	ID      int `json:"id"`
	Limit   int `json:"limite"`
	Balance int `json:"balance"`
}

type Statement struct {
	Balance_details   Balance        `json:"saldo"`
	Last_transactions []Transactions `json:"ultimas_transacoes"`
}

type Balance struct {
	Total          int       `json:"total"`
	Statement_date time.Time `json:"data_extrato"`
	Limit          int       `json:"limite"`
}

type TransactionResult struct {
	Limit   int `json:"limite"`
	Balance int `json:"saldo"`
}

// func (app Models) CreateClientModel(client Client) (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
// 	defer cancel()

// 	var newID int
// 	stmt := `insert into clientes (limite, saldo)
// 		values ($1, $2) returning id`

// 	err := db.QueryRowContext(ctx, stmt,
// 		client.Balance,
// 		client.Limit,
// 	).Scan(&newID)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return newID, nil
// }

func (app Models) GetTransactionsModel(clientId int) (*Statement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()

	var clientExists bool

	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM clientes WHERE id = $1) FOR UPDATE", clientId).Scan(&clientExists)
	if err != nil {
		return nil, err
	}

	if !clientExists {
		return nil, errors.New("cliente não encontrado")
	}

	// Consulta usando QueryRow para obter apenas uma linha
	row := tx.QueryRow(ctx, "SELECT saldo, limite, now() FROM clientes WHERE id = $1 FOR UPDATE", clientId)
	var balance Balance
	err = row.Scan(
		&balance.Total,
		&balance.Limit,
		&balance.Statement_date,
	)
	if err != nil {
		return nil, err
	}

	// Consulta usando Query para obter várias linhas
	rows, err := tx.Query(ctx, "SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id = $1 ORDER BY realizada_em DESC LIMIT 10 FOR UPDATE", clientId)
	if err != nil {
		log.Panicln(err)
		return nil, err
	}
	defer rows.Close()

	var transactions []Transactions

	for rows.Next() {
		var transaction Transactions
		err := rows.Scan(
			&transaction.Value,
			&transaction.Type,
			&transaction.Description,
			&transaction.Realizada_em,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	result := Statement{
		Balance_details:   balance,
		Last_transactions: transactions,
	}

	return &result, nil
}

func (app Models) CreateTransactionModel(transaction Transactions, clientId int) (*TransactionResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}
		err = tx.Commit(ctx)
	}()
	// Verificar se o cliente existe
	var clientExists bool
	err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM clientes WHERE id = $1) FOR UPDATE", clientId).Scan(&clientExists)
	if err != nil {
		return nil, err
	}
	if !clientExists {
		return nil, errors.New("cliente não encontrado")
	}

	// Verificar se o tipo de transação é válido
	if transaction.Type != "c" && transaction.Type != "d" {
		return nil, errors.New("tipo de transação inválido")
	}

	// Verificar se a descrição tem entre 1 e 10 caracteres
	if len(transaction.Description) < 1 || len(transaction.Description) > 10 {
		return nil, errors.New("descrição deve ter entre 1 e 10 caracteres")
	}

	rows := tx.QueryRow(ctx, "SELECT saldo, limite FROM clientes WHERE id = $1 FOR UPDATE", clientId)
	if err != nil {
		return nil, err
	}

	var balance Balance
	err = rows.Scan(
		&balance.Total,
		&balance.Limit,
	)
	if err != nil {
		return nil, err
	}

	var newBalance int

	if transaction.Type == "d" {
		newBalance = balance.Total - transaction.Value
		// Verificar se a transação de débito deixa o saldo inconsistente
		if newBalance < -balance.Limit {
			return nil, errors.New("a transação de débito deixaria o saldo inconsistente")
		}

		_, err = tx.Exec(ctx, "UPDATE clientes SET saldo = $1 WHERE id = $2", newBalance, clientId)
		if err != nil {
			return nil, err
		}
	}

	// Inserir a transação no banco de dados
	_, err = tx.Exec(ctx, "INSERT INTO transacoes (valor, tipo, descricao, cliente_id) VALUES ($1, $2, $3, $4)",
		transaction.Value, transaction.Type, transaction.Description, clientId)
	if err != nil {
		return nil, err
	}

	if transaction.Type == "c" {
		newBalance = balance.Total + transaction.Value
		_, err = tx.Exec(ctx, "UPDATE clientes SET saldo = $1 WHERE id = $2", newBalance, clientId)
		if err != nil {
			return nil, err
		}
	}

	result := &TransactionResult{
		Balance: newBalance,
		Limit:   balance.Limit,
	}

	return result, nil
}
