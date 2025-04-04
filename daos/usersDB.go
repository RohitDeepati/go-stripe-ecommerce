package daos

import (
	"database/sql"
	"fmt"
	"go-ecommerce/beans"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func InsertNewUser(db *sqlx.DB, newUser beans.Users)error{
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil{
		return err
	}
	newUser.Password = string(hashPassword)

	query := `Insert into users(name, email, password, role) values(:name, :email, :password, :role)`

	_, err = db.NamedExec(query, &newUser)
	if err != nil{
		return err
	}
	return nil
}

func GetUsers(db *sqlx.DB)([]beans.Users, error){
	var users []beans.Users
	query := `select userid, name, email, role from users`
	err := db.Select(&users, query)
	if err != nil{
		return nil, err
	}
	return users, err
}

func CheckingUserByEmailId(db *sqlx.DB, email string)(bool, error){
	var user beans.Users
	query := `select name, email, password, role from users where email = $1`
	err := db.Get(&user, query, email)
	if err != nil{
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}
	return true, nil
}

func CheckingUserLoginDetails(db *sqlx.DB, email string, password string)(bool, error){
	var hashPassword string
	query := `select password from users where email = $1`
	err := db.Get(&hashPassword, query, email)
	if err != nil{
		if err == sql.ErrNoRows{
			return false, nil
		}
		return false, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return false, nil
	}
	return true, nil
}

func QueryUserByEmail(db *sqlx.DB, email string)(beans.Users, error){
	var user beans.Users
	query := `select userid, name, email, password, role from users where email = $1`
	err := db.Get(&user, query, email)
	if err != nil{
		if err == sql.ErrNoRows{
			return user, fmt.Errorf("no user found with the given email: %v", email)
		}
		return user, err
	}
	return user, nil
}

func DeleteUserByEmail(db *sqlx.DB, email string)error{
	query := `delete from users where email = $1`
	res, err := db.Exec(query, email)
	
	if err != nil{
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil{
		return err
	}
	if rowsAffected == 0{
		return fmt.Errorf("no user found with the given email id")
	}
	return nil
}

func CheckingProductById(db *sqlx.DB, id string)(bool, error){
	var product beans.Products
	query := `SELECT name, title, price, stock, image from product WHERE productid=$1`
	err := db.Get(&product, query, id)
	if err != nil{
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}