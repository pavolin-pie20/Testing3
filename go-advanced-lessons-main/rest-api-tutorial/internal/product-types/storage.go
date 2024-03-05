package product_types

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, productType *ProductTypes) error
	FindAll(ctx context.Context) (u []ProductTypes, err error)
	FindOne(ctx context.Context, id string) (ProductTypes, error)
	Update(ctx context.Context, productType ProductTypes) error
	Delete(ctx context.Context, id string) error
	Backup(ctx context.Context, filePath string) error
}
