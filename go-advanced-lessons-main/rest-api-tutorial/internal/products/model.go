package products

import product_types "restapi-lesson/internal/product-types"

type Product struct {
	ProductID     string                     `json:"product_id"`
	TypeID        string                     `json:"type_id"`
	ProductName   string                     `json:"product_name"`
	Weight        float64                    `json:"weight"`
	Unit          string                     `json:"unit"`
	Description   string                     `json:"description"`
	PricePickup   float64                    `json:"price_pickup"`
	PriceDelivery float64                    `json:"price_delivery"`
	ProductType   product_types.ProductTypes `json:"product_type"`
}
