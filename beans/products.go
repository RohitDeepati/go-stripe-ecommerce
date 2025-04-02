package beans

type Products struct {
	ProductID int     `json:"productId" db:"productid"`
	SellerID  int     `json:"sellerId" db:"sellerid"`
	Name      string  `json:"name" db:"name"`
	Title     string  `json:"title" db:"title"`
	Price     float64 `json:"price" db:"price"`
	Stock     int     `json:"stock" db:"stock"`
	Image     string  `json:"image" db:"image"`
}