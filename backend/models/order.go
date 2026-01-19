package models

import (
	"time"
)

type Order struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Type        string    `json:"type"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateOrder(userID int64, orderType string, amount float64, description string) (*Order, error) {
	result, err := DB.Exec(
		`INSERT INTO orders (user_id, type, amount, description) VALUES (?, ?, ?, ?)`,
		userID, orderType, amount, description,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return GetOrderByID(id)
}

func GetOrderByID(id int64) (*Order, error) {
	o := &Order{}
	err := DB.QueryRow(
		`SELECT id, user_id, type, amount, description, created_at FROM orders WHERE id = ?`, id,
	).Scan(&o.ID, &o.UserID, &o.Type, &o.Amount, &o.Description, &o.CreatedAt)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func GetUserOrders(userID int64, page, pageSize int) ([]*Order, int, error) {
	var total int
	err := DB.QueryRow(`SELECT COUNT(*) FROM orders WHERE user_id = ?`, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		`SELECT id, user_id, type, amount, description, created_at
		 FROM orders WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		o := &Order{}
		err := rows.Scan(&o.ID, &o.UserID, &o.Type, &o.Amount, &o.Description, &o.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

func GetAllOrders(page, pageSize int) ([]*Order, int, error) {
	var total int
	err := DB.QueryRow(`SELECT COUNT(*) FROM orders`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		`SELECT id, user_id, type, amount, description, created_at
		 FROM orders ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*Order
	for rows.Next() {
		o := &Order{}
		err := rows.Scan(&o.ID, &o.UserID, &o.Type, &o.Amount, &o.Description, &o.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

type OrderWithUser struct {
	Order
	Username string `json:"username"`
}

func GetAllOrdersWithUser(page, pageSize int, orderType string) ([]*OrderWithUser, int, error) {
	var total int
	var countQuery string
	var queryArgs []interface{}

	if orderType != "" {
		countQuery = `SELECT COUNT(*) FROM orders WHERE type = ?`
		queryArgs = append(queryArgs, orderType)
	} else {
		countQuery = `SELECT COUNT(*) FROM orders`
	}

	err := DB.QueryRow(countQuery, queryArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var query string
	if orderType != "" {
		query = `SELECT o.id, o.user_id, o.type, o.amount, o.description, o.created_at, u.username
				 FROM orders o LEFT JOIN users u ON o.user_id = u.id
				 WHERE o.type = ?
				 ORDER BY o.created_at DESC LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, pageSize, offset)
	} else {
		query = `SELECT o.id, o.user_id, o.type, o.amount, o.description, o.created_at, u.username
				 FROM orders o LEFT JOIN users u ON o.user_id = u.id
				 ORDER BY o.created_at DESC LIMIT ? OFFSET ?`
		queryArgs = []interface{}{pageSize, offset}
	}

	rows, err := DB.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []*OrderWithUser
	for rows.Next() {
		o := &OrderWithUser{}
		err := rows.Scan(&o.ID, &o.UserID, &o.Type, &o.Amount, &o.Description, &o.CreatedAt, &o.Username)
		if err != nil {
			return nil, 0, err
		}
		orders = append(orders, o)
	}

	return orders, total, nil
}

// 统计相关
func GetTotalRecharge() (float64, error) {
	var total float64
	err := DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM orders WHERE type = 'recharge'`).Scan(&total)
	return total, err
}

func GetTotalConsume() (float64, error) {
	var total float64
	err := DB.QueryRow(`SELECT COALESCE(SUM(ABS(amount)), 0) FROM orders WHERE type = 'consume'`).Scan(&total)
	return total, err
}

func GetTodayRecharge() (float64, error) {
	var total float64
	err := DB.QueryRow(`SELECT COALESCE(SUM(amount), 0) FROM orders WHERE type = 'recharge' AND date(created_at) = date('now')`).Scan(&total)
	return total, err
}

func GetTodayConsume() (float64, error) {
	var total float64
	err := DB.QueryRow(`SELECT COALESCE(SUM(ABS(amount)), 0) FROM orders WHERE type = 'consume' AND date(created_at) = date('now')`).Scan(&total)
	return total, err
}
