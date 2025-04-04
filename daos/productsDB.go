package daos

import (
	"fmt"
	"go-ecommerce/beans"

	"github.com/jmoiron/sqlx"
)

func InsertNewProduct(db *sqlx.DB, newProduct beans.Products) error {
	var exists bool
	err := db.Get(&exists, "SELECT EXISTS (SELECT 1 FROM users WHERE userid = $1 AND role='Seller')", newProduct.SellerID)
	if err != nil{
		return err
	}
	if !exists{
		return fmt.Errorf("sellerid %d does not exist or is not a seller", newProduct.SellerID)
	}

	query := `INSERT INTO products (sellerid, name, title, price, stock, image) 
						VALUES (:sellerid, :name, :title, :price, :stock, :image)`
	_, err = db.NamedExec(query, &newProduct)
	if err != nil {
			return err
	}
	return nil
}

func GetProducts(db *sqlx.DB)([]beans.Products, error){
	var products []beans.Products
	query := `select productId,sellerid, name, title, price, stock, image from products`
	 err := db.Select(&products, query)
	 if err != nil{
		 return nil, err
	 }
	 return products, err
}

func RemoveProductById(db *sqlx.DB, productId int) error{
	query := `DELETE FROM PRODUCTS where productid=$1`

	result, err := db.Exec(query, productId)
	if err != nil{
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return err
	}
	if rowsAffected == 0{
		return fmt.Errorf("no product found with the id %d", productId)
	}
	return nil
}
