package types

type OrdersQueryParams struct {
	UserId    int    `json:"user_id" db:"user_id"`
	OrderKey  string `json:"order_key" db:"order_key" validate:"oneof=id price"`
	OrderType string `json:"order_type" db:"order_type" validate:"oneof=asc desc"`
	Page      int    `json:"page" db:"page" validate:"numeric,gt=0"`
	Limit     int    `json:"limit" db:"limit" validate:"numeric,gt=0"`
}

func GetDefaultOrdersQueryParams() OrdersQueryParams {
	return OrdersQueryParams{
		UserId:    0,
		OrderKey:  "id",
		OrderType: "asc",
		Page:      1,
		Limit:     5,
	}
}
