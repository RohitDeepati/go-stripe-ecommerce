package products

import (
	"go-ecommerce/beans"
	"go-ecommerce/daos"
	"go-ecommerce/middleware"
	"net/http"
	"strconv"

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
		authorized.DELETE("/products/:id", middleware.UserRoleMiddleware(), h.deleteProductById)
		authorized.PATCH("/products/:id", middleware.UserRoleMiddleware(), h.updateProductById)
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

func (h *ProductDB) deleteProductById(ctx *gin.Context){
	id := ctx.Param("id")
	productId, err := strconv.Atoi(id)

	if err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	err = daos.RemoveProductById(h.db, productId)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message":"Product deleted successfully"})
}

func (h *ProductDB) updateProductById(ctx *gin.Context){
	id := ctx.Param("id")

	var updateProduct beans.Products

	if err := ctx.ShouldBindJSON(&updateProduct); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	exists, err := daos.CheckingProductById(h.db, id)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":"failed to check the product with the given id", "details":err.Error()})
		return
	}
	if !exists{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"product is doesn't exist with the given id"})
		return
	}

	err = daos.UpdateProductById(h.db, id, &updateProduct)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message":"product updated"})
}