package products

type CreateProductsDTB struct {
	NameProduct          string  `json:"product_name"`
	TypeID               int     `json:"type_id"`
	WeightProduct        float64 `json:"weight"`
	UnitProduct          string  `json:"unit"`
	DescriptionProduct   string  `json:"description"`
	PricePickupProduct   float64 `json:"price_pickup"`
	PriceDeliveryProduct float64 `json:"price_delivery"`
}
