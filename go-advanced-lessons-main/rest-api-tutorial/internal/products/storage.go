package products

import "context"

type ProductRepository interface {
	CreateProduct(ctx context.Context, product *Product) error
	FindAllProduct(ctx context.Context) ([]Product, error)
	FindOneProduct(ctx context.Context, id string) (Product, error)
	UpdateProduct(ctx context.Context, product Product) error
	DeleteProduct(ctx context.Context, id string) error
	BackupProduct(ctx context.Context, filePath string) error
}
