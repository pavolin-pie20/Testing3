package customers

import (
	"context"
	"fmt"
	"restapi-lesson/internal/customers"
	"restapi-lesson/pkg/client/postgresql"
	"restapi-lesson/pkg/logging"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) CreateCustomer(ctx context.Context, customer *customers.Customer) error {
	q := `
		INSERT INTO customers (entity_type, contact_name, address, phone, user_priority, login, password, email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING customer_id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	err := r.client.QueryRow(
		ctx,
		q,
		customer.EntityType,
		customer.ContactName,
		customer.Address,
		customer.Phone,
		customer.UserPriority,
		customer.Login,
		customer.Password,
		customer.EMail).Scan(&customer.CustomerID)

	if err != nil {
		r.logger.Errorf("Error executing SQL query: %v", err)
		return err
	}
	r.logger.Infof("New row added with customer_id: %v", customer.CustomerID)
	return nil
}

func (r *repository) FindAllCustomers(ctx context.Context) (u []customers.Customer, err error) {
	q := `
		SELECT customer_id, entity_type, contact_name, address, phone, user_priority, login, password, email FROM public.customers;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	customer := make([]customers.Customer, 0)

	for rows.Next() {
		var cst customers.Customer

		err = rows.Scan(&cst.CustomerID, &cst.EntityType, &cst.ContactName, &cst.Address, cst.Phone, &cst.UserPriority, &cst.Login, &cst.Password, &cst.EMail)
		if err != nil {
			return nil, err
		}

		customer = append(customer, cst)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customer, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) *repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
