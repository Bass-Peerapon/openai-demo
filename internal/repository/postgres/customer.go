package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Bass-Peerapon/openai-demo/domain"
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
)

type CustomerRepository struct {
	Client *sqlx.DB
}

func NewCustomerRepository(client *sqlx.DB) *CustomerRepository {
	return &CustomerRepository{
		Client: client,
	}
}

func (r *CustomerRepository) GetCustomer(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	query := `SELECT * FROM customers WHERE id = $1`
	row := r.Client.QueryRowxContext(ctx, query, id)

	customer := &domain.Customer{}
	if err := row.StructScan(customer); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return customer, nil
}
