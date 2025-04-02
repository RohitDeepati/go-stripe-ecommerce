package products

import (
	"go-ecommerce/beans"
	"go-ecommerce/daos"
	"go-ecommerce/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type ProductDB struct{
	db *sqlx.DB
}

func ProductsHandler(db *sqlx.DB) ProductDB{
	return ProductDB{
		db : db,
	}
}

func (h *ProductDB) RouteGroup(r *gin.Engine){

	routeGroup := r.Group("/")
	authorized := routeGroup.Group("/", middleware.AuthenticationMiddleware())
	{
		authorized.GET("/products", middleware.UserRoleMiddleware(), h.getAllProducts)
	}
	routeGroup.POST("/newproducts", h.insertNewProduct )
	
}

func (h *ProductDB) insertNewProduct(ctx *gin.Context){
	var newProduct beans.Products
	if err := ctx.ShouldBindJSON(&newProduct); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
	}

	err := daos.InsertNewProduct(h.db, newProduct)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}



	ctx.JSON(http.StatusOK, gin.H{"message": "product added successfully"})
}

func (h *ProductDB) getAllProducts(ctx *gin.Context){
	res, err := daos.GetProducts(h.db)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}