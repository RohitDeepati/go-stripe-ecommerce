package users

import (
	"fmt"
	"go-ecommerce/beans"
	"go-ecommerce/daos"
	"go-ecommerce/middleware"
	"go-ecommerce/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type UserDB struct{
	db *sqlx.DB
}

func UsersHandler(db *sqlx.DB) UserDB{
	return UserDB{
		db : db,
	}
}

func (h *UserDB) RouteGroup(r *gin.Engine){
	routeGroup := r.Group("/")
	authorized := routeGroup.Group("/", middleware.AuthenticationMiddleware())
	{
		// authorized.GET("/users", middleware.UserRoleMiddleware(), h.getUsers)
		authorized.GET("/users", middleware.UserRoleMiddleware(), h.getUserByEmail)
		authorized.DELETE("/users", middleware.UserRoleMiddleware(), h.DeleteUserByEmail)
	}
	routeGroup.POST("/signup", h.signUpNewUser)
	routeGroup.POST("/login", h.loginuser)
	// routeGroup.GET("/users", h.getUserByEmail)
	
	// routeGroup.GET("/users", h.getAllUsers)
	
}


func (h *UserDB) signUpNewUser(ctx *gin.Context){
	var newUser beans.Users
	if err := ctx.ShouldBindJSON(&newUser); err != nil{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
		return
	}

	if newUser.Name == "" && newUser.Email == "" && newUser.Password == "" && newUser.Role == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"All fields are required"})
		return
	}

	if newUser.Name == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"Name is required"})
		return
	}

	if newUser.Email == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"Email is required"})
		return
	}
	if newUser.Password == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"Password is required"})
		return
	}

	if newUser.Role == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"Role is required"})
		return
	}
	
	if !utils.IsEmailIsValid(newUser.Email){
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"Invalid email format"})
		return
	}

	userExists, err := daos.CheckingUserByEmailId(h.db, newUser.Email)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	if userExists{
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
		return
	}
	
	err = daos.InsertNewUser(h.db, newUser)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()})
		return
	}

	token, err := utils.GenerateJwtToken(newUser.Email, newUser.Role)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user created successfully", "token": token})
}

func (h *UserDB) loginuser(ctx *gin.Context){
	var loginuser beans.Loginuser
	if err := ctx.ShouldBindJSON(&loginuser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request data", "details": err.Error()})
		return
}

// Log the request data to check
fmt.Printf("Received login data: %v", loginuser)


	// checking user exists
	userExists, err := daos.CheckingUserLoginDetails(h.db, loginuser.Email, loginuser.Password)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !userExists{
		ctx.JSON(http.StatusUnauthorized, gin.H{"error":"invalid email or password"})
		return
	}

	// fetch the user to the role
	user, err := daos.QueryUserByEmail(h.db, loginuser.Email)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	// generate token based on the rol
	token, err := utils.GenerateJwtToken(loginuser.Email, user.Role)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"token":token, "userid" : user.Email})
}

// func (h *UserDB) getUsers(ctx *gin.Context){
// 	res, err := daos.GetUsers(h.db)
// 	if err != nil{
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, res)
// }

func (h *UserDB) getUserByEmail(ctx *gin.Context){
	email, _ := ctx.GetQuery("email")
	
	fmt.Println(email)
	if email == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"email parameter is required"})
		return
	}
	
	var user beans.Users
	user.Email = email
	res, err := daos.QueryUserByEmail(h.db, user.Email)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, res)
}

func (h *UserDB) DeleteUserByEmail(ctx *gin.Context){
	email, _ := ctx.GetQuery("email")

	if email == ""{
		ctx.JSON(http.StatusBadRequest, gin.H{"error":"email parameter is required"})
		return
	}
	err := daos.DeleteUserByEmail(h.db, email)
	if err != nil{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

