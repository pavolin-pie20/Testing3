package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"restapi-lesson/internal/config"
	"restapi-lesson/pkg/utils"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) //Выполняет запрос, который не возвращает строки
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)              //Выполняет запрос, который возвращает набор результатов
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row                     //Выполняет запрос,который ожидается вернутьодну строку, и возвращает объект pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)                                                 //Начинает новую транзакцию и возвращает объект pgx.Tx
}

func NewClient(ctx context.Context, maxAttempts int, sc config.StorageConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	err = repeatable.DoWithTries(func() error { //созданная папка
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second) //5 секунд
		defer cancel()

		pool, err = pgxpool.Connect(ctx, dsn) //pgxpoll - безопасное соедиение для pgx
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}

// 	^ Для повторного выполения попыток установки соединения, т.к. докер постгерса может не подключится с первого раза
//	| в repeatable.go хранится код с повторным запуском
//	| Если установка соединения не удается, программа завершается с фатальной ошибкой
