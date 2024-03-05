package product_types

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"os"
	"restapi-lesson/internal/product-types"
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

func (r *repository) Create(ctx context.Context, productType *product_types.ProductTypes) error {
	q := `
		INSERT INTO product_types (type_name) VALUES ($1) RETURNING type_id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	err := r.client.QueryRow(ctx, q, productType.NameType).Scan(&productType.IDType)
	if err != nil {
		r.logger.Errorf("Error executing SQL query: %v", err)
		return err
	}

	r.logger.Infof("New row added with type_id: %v", productType.IDType)
	return nil
}

func (r *repository) FindAll(ctx context.Context) (u []product_types.ProductTypes, err error) {
	q := `
		SELECT type_id, type_name FROM public.product_types;
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	product_type := make([]product_types.ProductTypes, 0)

	for rows.Next() {
		var prd product_types.ProductTypes

		err = rows.Scan(&prd.IDType, &prd.NameType)
		if err != nil {
			return nil, err
		}

		product_type = append(product_type, prd)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return product_type, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (product_types.ProductTypes, error) {
	q := `
		SELECT type_id, type_name FROM public.product_types WHERE type_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var ath product_types.ProductTypes
	err := r.client.QueryRow(ctx, q, id).Scan(&ath.IDType, &ath.IDType)
	if err != nil {
		return product_types.ProductTypes{}, err
	}

	return ath, nil
}

func (r *repository) Update(ctx context.Context, productType product_types.ProductTypes) error {
	q := `
        UPDATE product_types
        SET type_name = $2
        WHERE type_id = $1
    `
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	_, err := r.client.Exec(ctx, q, productType.IDType, productType.NameType)
	if err != nil {
		// Обработка ошибок
		return err
	}

	return nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	q := `
		DELETE FROM product_types
		WHERE type_id = $1
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

func (r *repository) Backup(ctx context.Context, filePath string) error {
	q := `SELECT type_id, type_name FROM public.product_types`

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
	err = writer.Write([]string{"type_id", "type_name"})
	if err != nil {
		return err
	}

	// Записываем данные
	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		err = writer.Write([]string{id, name})
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

func NewRepository(client postgresql.Client, logger *logging.Logger) product_types.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
