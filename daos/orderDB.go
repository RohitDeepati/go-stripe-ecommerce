package daos

import (
	"fmt"
	"go-ecommerce/beans"

	"github.com/jmoiron/sqlx"
)

// func InsertItemToCart(db *sqlx.DB, products []beans.Orders)error{
// 	query := `insert into orders(userid, productid, quantity) values(:userid, :productid, :quantity)`
// 	_,err := db.NamedExec(query, products)
// 	if err != nil{
// 		return err
// 	}
// 	return nil
// }

func InsertItemToOrder(db *sqlx.DB, products []beans.OrderItem, userId int)error{
	tx, err := db.Beginx()
	if err != nil{
		return err
	}
	defer tx.Rollback()

	var orderID int
	err = tx.QueryRow("INSERT INTO orders (userid) VALUES ($1) RETURNING id", userId).Scan(&orderID)
	if err != nil{
		return err
	}

	query := `INSERT INTO order_items(order_id, userid, productid, quantity) VALUES(:order_id, :userid, :productid, :quantity)`
	for i := range products {
		products[i].OrderID = orderID // Assign the new order ID

		available, err := CheckStock(tx, products[i].ProductID, products[i].Quantity)
		if err != nil{
			tx.Rollback()
			return err
		}
		if !available{
			tx.Rollback()
			return fmt.Errorf("Not enough stock available for product Id %d", products[i].ProductID)
		}

		_, err = tx.NamedExec(query, products[i]) // Use individual product
		if err != nil {
			tx.Rollback()
			return err
		}

		err = DeductStock(tx, products[i].ProductID, products[i].Quantity)
		if err != nil{
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}


func QueryAllItems(db *sqlx.DB)([]beans.OrderResponse, error){
	var items []beans.OrderResponse
	query := `SELECT
	orders.id AS order_id,
	users.userid AS user_id,
	users.name AS user_name,
	products.productid AS product_id,
	products.name AS product_name,
	products.price AS product_price,
	orders.quantity AS order_quantity
FROM
	orders
JOIN
	users ON orders.userid = users.userId
JOIN
	products ON orders.productid = products.productId `
	err := db.Select(&items, query)
	if err != nil{
		return nil, err
	}
	return items, err
}

func QueryOrderByEmail(db *sqlx.DB, email string)([]beans.OrderResponse, error){
	var orderItems []beans.OrderResponse
	query := `SELECT
		order_items.id AS id,
    order_items.order_id AS order_item_id,
    users.userid AS user_id,
    users.name AS user_name,
    users.email AS email,
    products.productId AS product_id,
    products.image AS image,
    products.name AS product_name,
    products.title AS product_title,
    products.price AS product_price,
    order_items.quantity AS order_quantity
FROM
    orders
JOIN
    order_items ON orders.id = order_items.order_id  
JOIN
    users ON orders.userid = users.userId  
JOIN
    products ON order_items.productid = products.productId
	where 
		users.email = $1`
	err := db.Select(&orderItems, query, email)
	if err != nil{
		return nil, err
	}
	return orderItems, err
}

func CheckStock(exec sqlx.Queryer, productId int, quantity int) (bool, error){
	var stock int 
	err := exec.QueryRowx(`Select stock from products where productid=$1`, productId).Scan(&stock)
	if err != nil{
		return false, err
	}
	return stock >= quantity, nil
}

func DeductStock(exec sqlx.Execer, productId int, quantity int)error{
	result, err := exec.Exec(
		`Update products
		set stock = stock - $1
		where productId = $2 and stock >= $1`, quantity, productId)
		if err != nil{
			return nil
		}	

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0{
			return fmt.Errorf("not enough stock available")
		}
		return nil
}