package models

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	// 确保目录存在
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_foreign_keys=on")
	if err != nil {
		return err
	}

	// 测试连接
	if err = DB.Ping(); err != nil {
		return err
	}

	// 创建表
	if err = createTables(); err != nil {
		return err
	}

	log.Println("数据库初始化完成")
	return nil
}

func createTables() error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(100),
			balance DECIMAL(10,2) DEFAULT 0,
			max_ports INTEGER DEFAULT 3,
			max_devices INTEGER DEFAULT 10,
			role VARCHAR(20) DEFAULT 'user',
			status VARCHAR(20) DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS ports (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			port INTEGER UNIQUE NOT NULL,
			protocol VARCHAR(20) DEFAULT 'autojs',
			name VARCHAR(100) DEFAULT '',
			token VARCHAR(64) NOT NULL,
			status VARCHAR(20) DEFAULT 'active',
			expire_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS devices (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			port_id INTEGER NOT NULL,
			device_id VARCHAR(100) NOT NULL,
			device_name VARCHAR(100),
			device_info TEXT,
			last_ping INTEGER DEFAULT 0,
			status VARCHAR(20) DEFAULT 'offline',
			connected_at DATETIME,
			disconnected_at DATETIME,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (port_id) REFERENCES ports(id)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			type VARCHAR(20) NOT NULL,
			amount DECIMAL(10,2) NOT NULL,
			description VARCHAR(255),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS configs (
			key VARCHAR(50) PRIMARY KEY,
			value TEXT NOT NULL,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS device_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			device_id INTEGER NOT NULL,
			level VARCHAR(20) DEFAULT 'info',
			message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (device_id) REFERENCES devices(id)
		)`,
		// 索引
		`CREATE INDEX IF NOT EXISTS idx_ports_user_id ON ports(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_ports_status ON ports(status)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_port_id ON devices(port_id)`,
		`CREATE INDEX IF NOT EXISTS idx_devices_status ON devices(status)`,
		`CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_device_logs_device_id ON device_logs(device_id)`,
	}

	for _, table := range tables {
		if _, err := DB.Exec(table); err != nil {
			return err
		}
	}

	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
