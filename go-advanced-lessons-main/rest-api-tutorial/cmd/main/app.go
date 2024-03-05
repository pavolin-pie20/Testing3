package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"restapi-lesson/internal/config"
	customers2 "restapi-lesson/internal/customers"
	"restapi-lesson/internal/customers/db"
	product_types2 "restapi-lesson/internal/product-types"
	"restapi-lesson/internal/product-types/db"
	products2 "restapi-lesson/internal/products"
	"restapi-lesson/internal/products/db"
	"restapi-lesson/pkg/client/postgresql"
	"restapi-lesson/pkg/logging"
	"strings"
	"time"
)

// DataRow представляет строктуру наборе результатов базы данных
type DataRow struct {
	TypeID   string
	TypeName string
}

func main() {
	// Получаем путь к текущему каталогу
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// ----------------------------
	// СБОР ПУТЕЙ К ФАЙЛАМ
	indexHTMLPath := filepath.Join(currentDir, "cmd/main/index.html") //index.html
	styleCSSPath := filepath.Join(currentDir, "cmd/main/style.css")   // style.css

	productsHTMLPath := filepath.Join(currentDir, "cmd/main/products.html") // products.html
	// ----------------------------

	// ----------------------------
	// КОНФИГУРАЦИЯ БАЗЫ ДАННЫХ
	cfg := config.GetConfig()
	// ----------------------------

	// ----------------------------
	// СОЗДАНИЕ КЛИЕНТА POSTGRESQL
	yourPostgresClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		log.Fatalf("Ошибка создания PostgreSQL клиента: %v", err)
	}
	yourLogger := logging.GetLogger()
	// ----------------------------

	// ----------------------------
	// СОЗДАНИЕ РЕПОЗИТОРИЯ
	repo := product_types.NewRepository(yourPostgresClient, yourLogger)
	repoProduct := products.NewRepository(yourPostgresClient, yourLogger)
	repoCustomer := customers.NewRepository(yourPostgresClient, yourLogger)
	// ----------------------------

	//----------------------------
	// РАСПОЗНОВАНИЕ HTML И CSS ФУНКЦИЯ
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/" {
			dataRows, err := repo.FindAll(context.TODO())
			if err != nil {
				http.Error(res, fmt.Sprintf("Запрос не выполнен: %v", err), http.StatusInternalServerError) //Ошибка 500
				return
			}

			// Parsing HTML template Парсирование HTML шаблона
			tmpl, err := template.ParseFiles(indexHTMLPath)
			if err != nil {
				http.Error(res, fmt.Sprintf("Не удалось парсирование шаблона: %v", err), http.StatusInternalServerError) // Ошибка 500
				return
			}

			// Выполнение шаблона с помощью dataRows
			err = tmpl.Execute(res, struct{ Rows []product_types2.ProductTypes }{dataRows})
			if err != nil {
				http.Error(res, fmt.Sprintf("Не удалось выполнить шаблон: %v", err), http.StatusInternalServerError) // Ошибка 500
			}

		} else if req.URL.Path == "/style.css" {
			log.Printf("Обуслуживание CSS-файла: %s\n", styleCSSPath)
			http.ServeFile(res, req, styleCSSPath)

		} else if req.URL.Path == "/products.html" {
			log.Printf("Обуслуживание HTML-файла: %s\n", productsHTMLPath)

			dataRows, err := repoProduct.FindAllProduct(context.TODO()) // Используйте функцию для получения продуктов
			if err != nil {
				http.Error(res, fmt.Sprintf("Запрос не выполнен: %v", err), http.StatusInternalServerError)
				return
			}

			dataRows1, err := repo.FindAll(context.TODO()) // Используйте функцию для получения типов продуктов
			if err != nil {
				http.Error(res, fmt.Sprintf("Запрос не выполнен: %v", err), http.StatusInternalServerError)
				return
			}

			tmpl, err := template.ParseFiles(productsHTMLPath)
			if err != nil {
				http.Error(res, fmt.Sprintf("Не удалось парсирование шаблона: %v", err), http.StatusInternalServerError)
				return
			}

			Rows := struct {
				Products     []products2.Product
				ProductTypes []product_types2.ProductTypes
			}{
				Products:     dataRows,
				ProductTypes: dataRows1,
			}

			err = tmpl.Execute(res, Rows)
			if err != nil {
				http.Error(res, fmt.Sprintf("Не удалось выполнить шаблон: %v", err), http.StatusInternalServerError)
			}
		}
	})
	// ----------------------------

	// ----------------------------
	// ДОБАВЛЕНИЕ ФУНКЦИЯ
	http.HandleFunc("/add", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Метод не разрешен", http.StatusMethodNotAllowed) // Ошибка 405
			return
		}

		var request struct {
			TypeName string `json:"TypeName"`
		}

		// Вывод ошибки при неправильном декодировании JSON
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			http.Error(res, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest) // Ошибка 400
			return
		}

		// Создание нового экземпляра ProductTypes с указанием NameType
		newType := product_types2.ProductTypes{NameType: request.TypeName}

		// Вызов метода Create в репозитории, чтобы добавить новую строку в таблицу БД
		err = repo.Create(context.TODO(), &newType)
		if err != nil {
			http.Error(res, fmt.Sprintf("Ошибка при создании новой строки: %v", err), http.StatusInternalServerError) // Ошибка 500
			return
		}

		// Ответ добавления данных в новой строке
		response := struct {
			IDType   string `json:"IDType"`
			NameType string `json:"NameType"`
		}{
			IDType:   newType.IDType,
			NameType: newType.NameType,
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(response)
	})
	// ----------------------------

	// ----------------------------
	// РЕГИСТРАЦИЯ НОВОГО ПОЛЬЗОВАТЕЛЯ
	// ----------------------------
	http.HandleFunc("/register", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Метод не разрешен", http.StatusMethodNotAllowed) // Ошибка 405
			return
		}

		var request struct {
			EntityType   string `json:"EntityType"`
			ContactName  string `json:"ContactName"`
			Address      string `json:"Address"`
			Phone        string `json:"Phone"`
			UserPriority string `json:"UserPriority"`
			Login        string `json:"Login"`
			Password     string `json:"Password"`
			EMail        string `json:"EMail"`
		}

		// Вывод ошибки при неправильном декодировании JSON
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			http.Error(res, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest) // Ошибка 400
			return
		}

		// Создание нового экземпляра ProductTypes с указанием NameType
		newCustomer := customers2.Customer{
			EntityType:   request.EntityType,
			ContactName:  request.ContactName,
			Address:      request.Address,
			Phone:        request.Phone,
			UserPriority: request.UserPriority,
			Login:        request.Login,
			Password:     request.Password,
			EMail:        request.EMail,
		}

		// Вызов метода Create в репозитории, чтобы добавить новую строку в таблицу БД
		err = repoCustomer.CreateCustomer(context.TODO(), &newCustomer)
		if err != nil {
			http.Error(res, fmt.Sprintf("Ошибка при создании новой строки: %v", err), http.StatusInternalServerError) // Ошибка 500
			return
		}

		// Ответ добавления данных в новой строке
		response := struct {
			CustomerID   string `json:"CustomerID"`
			EntityType   string `json:"EntityType"`
			ContactName  string `json:"ContactName"`
			Address      string `json:"Address"`
			Phone        string `json:"Phone"`
			UserPriority string `json:"UserPriority"`
			Login        string `json:"Login"`
			Password     string `json:"Password"`
			EMail        string `json:"EMail"`
		}{
			CustomerID:   newCustomer.CustomerID,
			EntityType:   newCustomer.EntityType,
			ContactName:  newCustomer.ContactName,
			Address:      newCustomer.Address,
			Phone:        newCustomer.Phone,
			UserPriority: newCustomer.UserPriority,
			Login:        newCustomer.Login,
			Password:     newCustomer.Password,
			EMail:        newCustomer.EMail,
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(response)
	})

	// products PRODUCTS products
	// Добавление фунции обработки HTTP-запроса для создания новой строки в products
	http.HandleFunc("/add_product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(res, "Метод не разрешен", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			TypeID        string  `json:"TypeID"`
			ProductName   string  `json:"ProductName"`
			Weight        float64 `json:"Weight"`
			Unit          string  `json:"Unit"`
			Description   string  `json:"Description"`
			PricePickup   float64 `json:"PricePickup"`
			PriceDelivery float64 `json:"PriceDelivery"`
		}

		// Вывод ошибки при неправильном декодировании JSON
		err := json.NewDecoder(req.Body).Decode(&request)
		if err != nil {
			http.Error(res, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest)
			return
		}

		// Создание нового экземпляра Product
		newProduct := products2.Product{
			TypeID:        request.TypeID,
			ProductName:   request.ProductName,
			Weight:        request.Weight,
			Unit:          request.Unit,
			Description:   request.Description,
			PricePickup:   request.PricePickup,
			PriceDelivery: request.PriceDelivery,
		}

		// Вызов метода CreateProduct в репозитории, чтобы добавить новую строку в таблицу БД
		err = repoProduct.CreateProduct(context.TODO(), &newProduct)
		if err != nil {
			http.Error(res, fmt.Sprintf("Ошибка при создании новой строки: %v", err), http.StatusInternalServerError)
			return
		}

		// Ответ добавления данных в новой строке
		response := struct {
			ProductID     string  `json:"ProductID"`
			TypeID        string  `json:"TypeID"`
			ProductName   string  `json:"ProductName"`
			Weight        float64 `json:"Weight"`
			Unit          string  `json:"Unit"`
			Description   string  `json:"Description"`
			PricePickup   float64 `json:"PricePickup"`
			PriceDelivery float64 `json:"PriceDelivery"`
		}{
			ProductID:     newProduct.ProductID,
			TypeID:        newProduct.TypeID,
			ProductName:   newProduct.ProductName,
			Weight:        newProduct.Weight,
			Unit:          newProduct.Unit,
			Description:   newProduct.Description,
			PricePickup:   newProduct.PricePickup,
			PriceDelivery: newProduct.PriceDelivery,
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(response)
	})

	// ----------------------------

	// ----------------------------
	// УДАЛЕНИЕ ФУНКЦИЯ
	http.HandleFunc("/delete/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			id := strings.TrimPrefix(req.URL.Path, "/delete/")
			if id != "" {
				err := repo.Delete(context.TODO(), id)
				if err != nil {
					http.Error(res, fmt.Sprintf("Не удалось удалить строку: %v", err), http.StatusInternalServerError) //Ошибка 500
					return
				}
				res.WriteHeader(http.StatusOK) // Статус 200
				return
			}
		}
		http.NotFound(res, req)
	})
	// ----------------------------

	// ----------------------------
	// УДАЛЕНИЕ ДЛЯ PRODUCT
	// УДАЛЕНИЕ ФУНКЦИЯ
	http.HandleFunc("/delete_product/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {
			id := strings.TrimPrefix(req.URL.Path, "/delete_product/")
			if id != "" {
				err := repoProduct.DeleteProduct(context.TODO(), id)
				if err != nil {
					http.Error(res, fmt.Sprintf("Не удалось удалить строку: %v", err), http.StatusInternalServerError) //Ошибка 500
					return
				}
				res.WriteHeader(http.StatusOK) // Статус 200
				return
			}
		}
		http.NotFound(res, req)
	})
	// ----------------------------

	// ----------------------------
	// РЕДАКТОР ФУНКЦИЯ
	http.HandleFunc("/edit/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {

			var request struct {
				IDType   string `json:"IDType"`
				TypeName string `json:"TypeName"`
			}

			// Вывод ошибки при неправильном декодировании JSON
			err := json.NewDecoder(req.Body).Decode(&request)
			if err != nil {
				http.Error(res, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest) // Ошибка 400
				return
			}

			// Печать данных для проверки их корректности
			fmt.Printf("Полученные IDType: %s, TypeName: %s\n", request.IDType, request.TypeName)

			// Вызов метода Update в репозитории, чтобы редактировать строку в таблице БД
			err = repo.Update(context.TODO(), product_types2.ProductTypes{
				IDType:   request.IDType,
				NameType: request.TypeName,
			})

			//Вывод ошибки, в случае неуспешной попытки редактирования строки
			if err != nil {
				http.Error(res, fmt.Sprintf("Ошибка при попытке редактирования строки: %v", err), http.StatusInternalServerError) // Ошибка 500
				return
			}

			res.WriteHeader(http.StatusOK) // Статус 200
			return
		}
		http.NotFound(res, req)
	})
	// ----------------------------

	// ----------------------------
	// РЕДАКТОР ДЛЯ PRODUCTS
	// РЕДАКТОР ФУНКЦИЯ
	http.HandleFunc("/edit_product/", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {

			var request struct {
				ProductID     string  `json:"ProductID"`
				TypeID        string  `json:"TypeID"`
				ProductName   string  `json:"ProductName"`
				Weight        float64 `json:"Weight"`
				Unit          string  `json:"Unit"`
				Description   string  `json:"Description"`
				PricePickup   float64 `json:"PricePickup"`
				PriceDelivery float64 `json:"PriceDelivery"`
			}

			// Вывод ошибки при неправильном декодировании JSON
			err := json.NewDecoder(req.Body).Decode(&request)
			if err != nil {
				http.Error(res, fmt.Sprintf("Ошибка декодирования JSON: %v", err), http.StatusBadRequest) // Ошибка 400
				return
			}

			// Печать данных для проверки их корректности
			fmt.Printf("Полученные ProductID: #{request.ProductID}, TypeID: #{request.TypeID}, ProductName: #{request.ProductName}, Weight: #{request.Weight}, Unit: #{request.Unit}, Description: #{request.Description}, PricePickup: #{request.PricePickup}, PriceDelivery: #{request.PriceDelivery}")

			// Вызов метода Update в репозитории, чтобы редактировать строку в таблице БД
			err = repoProduct.UpdateProduct(context.TODO(), products2.Product{
				ProductID:     request.ProductID,
				TypeID:        request.TypeID,
				ProductName:   request.ProductName,
				Weight:        request.Weight,
				Unit:          request.Unit,
				Description:   request.Description,
				PricePickup:   request.PricePickup,
				PriceDelivery: request.PriceDelivery,
			})

			//Вывод ошибки, в случае неуспешной попытки редактирования строки
			if err != nil {
				http.Error(res, fmt.Sprintf("Ошибка при попытке редактирования строки: %v", err), http.StatusInternalServerError) // Ошибка 500
				return
			}

			res.WriteHeader(http.StatusOK) // Статус 200
			return
		}
		http.NotFound(res, req)
	})
	// ----------------------------

	// ----------------------------
	// БЭКАП ФУНКЦИЯ
	http.HandleFunc("/backup", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {

			// Путь, куда сохранить бэкап
			backupFilePath := "backup_" + time.Now().Format("20060102_150405") + ".csv"

			// Вызов функциИ Backup, чтобы создать бэкап
			err := repo.Backup(context.TODO(), backupFilePath)
			if err != nil {
				http.Error(res, fmt.Sprintf("Backup failed: %v", err), http.StatusInternalServerError) // Ошибка 500
				return
			}

			// Ответ, что бэкап успешно создан
			res.WriteHeader(http.StatusOK) // Статус 200
			return
		}
		http.NotFound(res, req)
	})
	// ----------------------------

	// ----------------------------
	// БЭКАП PRODUCT
	// БЭКАП ФУНКЦИЯ
	http.HandleFunc("/backup_product", func(res http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPost {

			// Путь, куда сохранить бэкап
			backupFilePath := "backup_product" + time.Now().Format("20060102_150405") + ".csv"

			// Вызов функциИ Backup, чтобы создать бэкап
			err := repoProduct.BackupProduct(context.TODO(), backupFilePath)
			if err != nil {
				http.Error(res, fmt.Sprintf("Backup failed: %v", err), http.StatusInternalServerError) // Ошибка 500
				return
			}

			// Ответ, что бэкап успешно создан
			res.WriteHeader(http.StatusOK) // Статус 200
			return
		}
		http.NotFound(res, req)
	})
	// ----------------------------

	// ----------------------------
	// ЗАПУСК СЕРВЕРА НА ПОРТУ "1234" С TLS
	address := fmt.Sprintf(":%s", cfg.Listen.Port)
	log.Printf("Сервер запущен на порту %s", address)

	// localhost.crt   |   СЕРТИФИКАТ  |   для запуска сервера на HTTPS
	// localhost.key   |      КЛЮЧ     |   для запуска сервара на HTTPS

	log.Fatal(http.ListenAndServeTLS(address, "localhost.crt", "localhost.key", nil))
	// ----------------------------

}
