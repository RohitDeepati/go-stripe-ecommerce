package beans

type Cart struct{
	Name     string `json:"productName"`
	Title 	string	`json:"productTitle"`
	Quantity int64  `json:"quantity"`
	Price    int64  `json:"price"`    
	ImageURL string `json:"image_url"` 
}