package repository

import (
	"2020_1_drop_table/internal/app/customer"
	"2020_1_drop_table/internal/app/customer/models"
	"context"
	"github.com/jmoiron/sqlx"
)

type postgresCustomerRepository struct {
	conn *sqlx.DB
}

func NewPostgresCustomerRepository(conn *sqlx.DB) customer.Repository {
	return &postgresCustomerRepository{
		conn: conn,
	}
}

func (p *postgresCustomerRepository) Add(ctx context.Context, cu models.Customer) (models.Customer, error) {
	query := `INSERT INTO Customer(
    CafeID, 
	Points) 
	VALUES ($1,$2) 
	RETURNING *`

	var customerDB models.Customer
	err := p.conn.GetContext(ctx, &customerDB, query, cu.CafeID, cu.Points)

	return customerDB, err
}

func (p *postgresCustomerRepository) SetLoyaltyPoints(ctx context.Context, points int,
	customerID string) (models.Customer, error) {

	query := `UPDATE Customer SET Points=$1 WHERE CustomerID=$2 RETURNING *`

	var customerDB models.Customer
	err := p.conn.GetContext(ctx, &customerDB, query, points, customerID)

	return customerDB, err
}

func (p *postgresCustomerRepository) GetByID(ctx context.Context, customerID string) (models.Customer, error) {
	query := `SELECT * FROM Customer WHERE CustomerID=$1`

	var customerDB models.Customer
	err := p.conn.GetContext(ctx, &customerDB, query, customerID)

	return customerDB, err
}

func (p *postgresCustomerRepository) DeleteByID(ctx context.Context, customerID string) error {
	query := `DELETE FROM Customer WHERE CustomerID=$1`

	_, err := p.conn.ExecContext(ctx, query, customerID)

	return err
}
