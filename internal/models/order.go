package models

import (
	"context"
	"time"
	"fmt"

	"github.com/CicadaHymn/guitar-shop-api/internal/db"

	)
type Order struct {
	ID        int    `json:"id"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	ProductID int    `json:"product_id"`
	Status    string `json:"status"`
}

func CreateOrder(b *Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO orders (phone, address, product_id, status)
	VALUES ($1, $2, $3, $4)
	RETURNING id`

	err := db.Pool.QueryRow(ctx, query,b.Phone, b.Address, b.ProductID, b.Status).Scan(&b.ID)
	
	if err != nil {
		return fmt.Errorf("ошибка создания заказа: %v", err)
	}
	return nil
	
}

func GetOrders() ([]Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, phone, address, product_id, status
	FROM orders`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения заказов: %v", err)
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		if err := rows.Scan(&order.ID, &order.Phone, &order.Address, &order.ProductID, &order.Status); err != nil {
			return nil, fmt.Errorf("ошибка сканирования заказа: %v", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}