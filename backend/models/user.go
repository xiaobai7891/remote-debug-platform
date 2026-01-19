package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	Password   string    `json:"-"`
	Email      string    `json:"email"`
	Balance    float64   `json:"balance"`
	MaxPorts   int       `json:"max_ports"`
	MaxDevices int       `json:"max_devices"`
	Role       string    `json:"role"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func CreateUser(username, password, email string) (*User, error) {
	result, err := DB.Exec(
		`INSERT INTO users (username, password, email, balance, max_ports, max_devices, role, status)
		 VALUES (?, ?, ?, 0, 3, 10, 'user', 'active')`,
		username, password, email,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return GetUserByID(id)
}

func GetUserByID(id int64) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		`SELECT id, username, password, email, balance, max_ports, max_devices, role, status, created_at, updated_at
		 FROM users WHERE id = ?`, id,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance,
		&user.MaxPorts, &user.MaxDevices, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByUsername(username string) (*User, error) {
	user := &User{}
	err := DB.QueryRow(
		`SELECT id, username, password, email, balance, max_ports, max_devices, role, status, created_at, updated_at
		 FROM users WHERE username = ?`, username,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance,
		&user.MaxPorts, &user.MaxDevices, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *User) Update() error {
	_, err := DB.Exec(
		`UPDATE users SET email = ?, balance = ?, max_ports = ?, max_devices = ?,
		 role = ?, status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		u.Email, u.Balance, u.MaxPorts, u.MaxDevices, u.Role, u.Status, u.ID,
	)
	return err
}

func (u *User) UpdatePassword(newPassword string) error {
	_, err := DB.Exec(`UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		newPassword, u.ID)
	return err
}

func (u *User) AddBalance(amount float64, description string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 更新余额
	_, err = tx.Exec(`UPDATE users SET balance = balance + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		amount, u.ID)
	if err != nil {
		return err
	}

	// 记录订单
	orderType := "recharge"
	if amount < 0 {
		orderType = "consume"
	}
	_, err = tx.Exec(`INSERT INTO orders (user_id, type, amount, description) VALUES (?, ?, ?, ?)`,
		u.ID, orderType, amount, description)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetAllUsers(page, pageSize int) ([]*User, int, error) {
	var total int
	err := DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		`SELECT id, username, password, email, balance, max_ports, max_devices, role, status, created_at, updated_at
		 FROM users ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance,
			&user.MaxPorts, &user.MaxDevices, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func GetUserPortCount(userID int64) (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM ports WHERE user_id = ? AND status = 'active'`, userID).Scan(&count)
	return count, err
}

func GetUserDeviceCount(userID int64) (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM devices WHERE user_id = ? AND status = 'online'`, userID).Scan(&count)
	return count, err
}

func UserExists(username string) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM users WHERE username = ?`, username).Scan(&count)
	return count > 0, err
}

func EmailExists(email string) (bool, error) {
	if email == "" {
		return false, nil
	}
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM users WHERE email = ?`, email).Scan(&count)
	return count > 0, err
}

func SearchUsers(keyword string, page, pageSize int) ([]*User, int, error) {
	var total int
	searchTerm := "%" + keyword + "%"
	err := DB.QueryRow(`SELECT COUNT(*) FROM users WHERE username LIKE ? OR email LIKE ?`,
		searchTerm, searchTerm).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		`SELECT id, username, password, email, balance, max_ports, max_devices, role, status, created_at, updated_at
		 FROM users WHERE username LIKE ? OR email LIKE ? ORDER BY id DESC LIMIT ? OFFSET ?`,
		searchTerm, searchTerm, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Balance,
			&user.MaxPorts, &user.MaxDevices, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func GetUserStats(userID int64) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 端口数
	var portCount int
	DB.QueryRow(`SELECT COUNT(*) FROM ports WHERE user_id = ?`, userID).Scan(&portCount)
	stats["port_count"] = portCount

	// 活跃端口数
	var activePortCount int
	DB.QueryRow(`SELECT COUNT(*) FROM ports WHERE user_id = ? AND status = 'active'`, userID).Scan(&activePortCount)
	stats["active_port_count"] = activePortCount

	// 在线设备数
	var onlineDeviceCount int
	DB.QueryRow(`SELECT COUNT(*) FROM devices WHERE user_id = ? AND status = 'online'`, userID).Scan(&onlineDeviceCount)
	stats["online_device_count"] = onlineDeviceCount

	// 总消费
	var totalConsume sql.NullFloat64
	DB.QueryRow(`SELECT SUM(ABS(amount)) FROM orders WHERE user_id = ? AND type = 'consume'`, userID).Scan(&totalConsume)
	if totalConsume.Valid {
		stats["total_consume"] = totalConsume.Float64
	} else {
		stats["total_consume"] = 0
	}

	return stats, nil
}
