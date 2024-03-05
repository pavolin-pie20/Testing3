package customers

import "context"

type CustomerRepository interface {
	CreateCustomer(ctx context.Context, customer *Customer) error
	FindAllCustomers(ctx context.Context) ([]Customer, error)
}
