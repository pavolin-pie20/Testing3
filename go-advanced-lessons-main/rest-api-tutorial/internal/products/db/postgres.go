package products

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"os"
	"restapi-lesson/internal/products"
	"restapi-lesson/pkg/client/postgresql"
	"restapi-lesson/pkg/logging"
	"strconv"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) CreateProduct(ctx context.Context, product *products.Product) error {
	q := `
		INSERT INTO products (type_id, product_name, weight, unit, description, price_pickup, price_delivery) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) 
		RETURNING product_id 
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	err := r.client.QueryRow(ctx, q, product.TypeID, product.ProductName, product.Weight, product.Unit, product.Description, product.PricePickup, product.PriceDelivery).Scan(&product.ProductID)
	if err != nil {
		r.logger.Errorf("Error executing SQL query: %v", err)
		return err
	}

	r.logger.Infof("New row added with product_id: %v", product.ProductID)
	return nil
}

func (r *repository) FindAllProduct(ctx context.Context) (u []products.Product, err error) {
	q := `
		SELECT p.product_id, p.type_id, pt.type_name, p.product_name, p.weight, p.unit, p.description, p.price_pickup, p.price_delivery	 
		FROM public.products p
		JOIN public.product_types pt ON p.type_id = pt.type_id;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	product := make([]products.Product, 0)

	for rows.Next() {
		var prd products.Product

		err = rows.Scan(&prd.ProductID, &prd.ProductType.IDType, &prd.ProductType.NameType, &prd.ProductName, &prd.Weight, &prd.Unit, &prd.Description, &prd.PricePickup, &prd.PriceDelivery)
		if err != nil {
			return nil, err
		}

		product = append(product, prd)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return product, nil
}

func (r *repository) UpdateProduct(ctx context.Context, product products.Product) error {
	q := `
        UPDATE products
        SET 
            type_id = $2,
            product_name = $3,
            weight = $4,
            unit = $5,
            description = $6,
            price_pickup = $7,
            price_delivery = $8
        WHERE product_id = $1
    `
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	_, err := r.client.Exec(ctx, q, product.ProductID, product.ProductType.IDType, product.ProductName, product.Weight, product.Unit, product.Description, product.PricePickup, product.PriceDelivery)
	if err != nil {
		// Обработка ошибок
		return err
	}

	return nil
}

func (r *repository) DeleteProduct(ctx context.Context, id string) error {
	q := `
		DELETE FROM products
		WHERE product_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	_, err := r.client.Exec(ctx, q, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) BackupProduct(ctx context.Context, filePath string) error {
	q := `
		SELECT product_id, type_id, product_name, weight, unit, description, price_pickup, price_delivery 
		FROM public.products
		`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	// Выполняем запрос для получения данных
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Открываем файл для записи
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Создаем CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Записываем заголовок
	err = writer.Write([]string{"product_id", "type_id", "product_name", "weight", "unit", "description", "price_pickup", "price_delivery"})
	if err != nil {
		return err
	}

	// Записываем данные
	for rows.Next() {
		var productid, typeid, productname, unitproduct, descriptionproduct string
		var weightproduct, pricepickup, pricedelivery float64

		err := rows.Scan(&productid, &typeid, &productname, &weightproduct, &unitproduct, &descriptionproduct, &pricepickup, &pricedelivery)
		if err != nil {
			return err
		}

		// Преобразование float64 в строки
		weightStr := strconv.FormatFloat(weightproduct, 'f', -1, 64)
		pricepickupStr := strconv.FormatFloat(pricepickup, 'f', -1, 64)
		pricedeliveryStr := strconv.FormatFloat(pricedelivery, 'f', -1, 64)

		// Запись в CSV
		err = writer.Write([]string{productid, typeid, productname, weightStr, unitproduct, descriptionproduct, pricepickupStr, pricedeliveryStr})
		if err != nil {
			return err
		}
	}

	// Проверяем наличие ошибок при обработке строк
	if err = rows.Err(); err != nil {
		return err
	}

	return nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) *repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
