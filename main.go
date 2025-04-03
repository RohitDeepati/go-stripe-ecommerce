package main

import (
	"go-ecommerce/controllers/order"
	"go-ecommerce/controllers/products"
	"go-ecommerce/controllers/stripe"
	"go-ecommerce/controllers/users"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)


func initDB()(*sqlx.DB, error){
	connStr := "host=centerbeam.proxy.rlwy.net port=53668 user=postgres password=SQTxYdfQMgVqVYhAujhmyBsOLwMBTmXA dbname=railway sslmode=disable"
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil{
		return nil, err
	}
	return db, nil
}

func main(){
	db, err := initDB()

	if err != nil{
		panic(err)
	}

	defer db.Close()

	rg := gin.Default()
	rg.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // React app's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	userController := users.UsersHandler(db)
	productController := products.ProductsHandler(db)
	orderController := order.OrderHandler(db)
	paymentController := stripe.StripeHandler(db)

	userController.RouteGroup(rg)
	productController.RouteGroup(rg)
	orderController.RouteGroup(rg)
	paymentController.RouteGroup(rg)
	rg.Run(":9090")
}