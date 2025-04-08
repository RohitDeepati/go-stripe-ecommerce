package beans



type Order struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
}

type OrderItem struct {
	OrderID   int `json:"order_id" db:"order_id"`
	UserID    int `json:"userid" db:"user_id"`
	ProductID int `json:"productid" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity"`
}


type OrderResponse struct{
	ID			int		`json:"id" db:"id"`
	OrderID int `json:"order_id" db:"order_item_id"`
    UserID      int    `json:"user_id" db:"user_id"`
    UserName    string `json:"user_name" db:"user_name"`
		Email 			string	`json:"email" db:"email"`
		Image				string	`json:"image" db:"image"`
    ProductID   int    `json:"product_id" db:"product_id"`
    ProductName string `json:"product_name" db:"product_name"`
		ProductTitle string	`json:"product_title" db:"product_title"`
    ProductPrice float64   `json:"product_price" db:"product_price"`
    OrderQuantity int  `json:"order_quantity" db:"order_quantity"`
}