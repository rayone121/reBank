package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rayone121/reBank/backend/config"
	"github.com/rayone121/reBank/backend/types"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountByID(int) (*types.Account, error)
	GetAccountByNumber(int64) (*types.Account, error)
	GetAccountByUserName(string) (*types.Account, error)
}

type PostgressStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgressStore, error) {
	connStr := config.GetDBConnectionString()
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgressStore{
		db: db,
	}, nil
}

func (s *PostgressStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgressStore) createAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		username varchar(32), 
		first_name varchar(32),
		last_name varchar(32),
		encrypted_password varchar(64),
		number serial,
		balance integer,
		created_at timestamp
		)`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgressStore) Close() {
	s.db.Close()
}

func (s *PostgressStore) CreateAccount(a *types.Account) error {
	query := `INSERT INTO account 
	(username, first_name, last_name, encrypted_password, number, balance, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	resp, err := s.db.Query(
		query,
		a.UserName,
		a.FirstName,
		a.LastName,
		a.EncryptedPassword,
		a.Number,
		a.Balance,
		a.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("Response: %+v\n", resp)

	return nil
}

func (s *PostgressStore) UpdateAccount(a *types.Account) error {
	return nil
}

func (s *PostgressStore) DeleteAccount(id int) error {
	_, err := s.db.Query(`DELETE FROM account WHERE id = $1`, id)
	return err
}

func (s *PostgressStore) GetAccountByUserName(username string) (*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE username = $1`, username)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAcc(rows)
	}
	return nil, fmt.Errorf("account with username %s not found", username)
}

func (s *PostgressStore) GetAccountByNumber(number int64) (*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE number = $1`, number)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAcc(rows)
	}
	return nil, fmt.Errorf("account with number %d not found", number)
}

func (s *PostgressStore) GetAccountByID(id int) (*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoAcc(rows)
	}
	return nil, fmt.Errorf("account with id %d not found", id)
}

func (s *PostgressStore) GetAccounts() ([]*types.Account, error) {
	rows, err := s.db.Query(`SELECT * FROM account`)
	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}

	for rows.Next() {
		account, err := scanIntoAcc(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}

func scanIntoAcc(rows *sql.Rows) (*types.Account, error) {
	account := new(types.Account)
	err := rows.Scan(
		&account.ID,
		&account.UserName,
		&account.FirstName,
		&account.LastName,
		&account.EncryptedPassword,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	return account, err
}
