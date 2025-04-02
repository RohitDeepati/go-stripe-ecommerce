package order

import (
	"go-ecommerce/beans"
	"go-ecommerce/daos"
	"go-ecommerce/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type OrderDB struct{
	db *sqlx.DB
}

func OrderHandler(db *sqlx.DB) OrderDB{
	return OrderDB{
		db : db,
	}
}

func (h *OrderDB) RouteGroup(r *gin.Engine){

	routeGroup := r.Group("/")
	authorized := routeGroup.Group("/", middleware.AuthenticationMiddleware())
	{
		authorized.GET("/orders", middleware.UserRoleMiddleware(), h.getOrdersByEmail)
	}
	routeGroup.POST("/orders", h.insertItemtoOrder)
	// routeGroup.GET("/items", h.getItemsInCart)
}

func (h *OrderDB) insertItemtoOrder(ctx *gin.Context){
	var orderItems []beans.OrderItem

	if err := ctx.ShouldBindJSON(&orderItems); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	if len(orderItems) == 0{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Order must contain at least one item"})
		return
	}

	userID := orderItems[0].UserID

	err := daos.InsertItemToOrder(h.db, orderItems, userID)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Items successfully added to the cart"})
}


// func (h *OrderDB) getItemsAllItems(ctx *gin.Context){
// 	res, err := daos.QueryAllItems(h.db)
// 	if err != nil{
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, res)
// }



func (h *OrderDB) getOrdersByEmail(ctx *gin.Context){
	email, _ := ctx.GetQuery("email")
	if email == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"required email id"})
		return
	}

	exists, err := daos.CheckingUserByEmailId(h.db, email)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	if !exists{
		ctx.JSON(http.StatusNotFound, gin.H{"error":"no account found with the given email address"})
		return
	}
	
	res, err := daos.QueryOrderByEmail(h.db, email)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)	
}

