package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

type Device struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"user_id"`
	PortID         int64      `json:"port_id"`
	DeviceID       string     `json:"device_id"`
	DeviceName     string     `json:"device_name"`
	DeviceInfo     string     `json:"device_info"`
	LastPing       int        `json:"last_ping"`
	Status         string     `json:"status"`
	ConnectedAt    *time.Time `json:"connected_at"`
	DisconnectedAt *time.Time `json:"disconnected_at"`
}

type DeviceInfo struct {
	Brand   string `json:"brand"`
	Model   string `json:"model"`
	OS      string `json:"os"`
	Version string `json:"version"`
	SDK     int    `json:"sdk"`
}

func CreateOrUpdateDevice(userID, portID int64, deviceID, deviceName string, info *DeviceInfo) (*Device, error) {
	infoJSON, _ := json.Marshal(info)

	// 检查设备是否存在
	var existingID int64
	err := DB.QueryRow(`SELECT id FROM devices WHERE port_id = ? AND device_id = ?`, portID, deviceID).Scan(&existingID)

	if err == sql.ErrNoRows {
		// 创建新设备
		result, err := DB.Exec(
			`INSERT INTO devices (user_id, port_id, device_id, device_name, device_info, status, connected_at)
			 VALUES (?, ?, ?, ?, ?, 'online', CURRENT_TIMESTAMP)`,
			userID, portID, deviceID, deviceName, string(infoJSON),
		)
		if err != nil {
			return nil, err
		}
		id, _ := result.LastInsertId()
		return GetDeviceByID(id)
	}

	// 更新已存在的设备
	_, err = DB.Exec(
		`UPDATE devices SET device_name = ?, device_info = ?, status = 'online', connected_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		deviceName, string(infoJSON), existingID,
	)
	if err != nil {
		return nil, err
	}

	return GetDeviceByID(existingID)
}

func GetDeviceByID(id int64) (*Device, error) {
	d := &Device{}
	var connectedAt, disconnectedAt sql.NullTime
	err := DB.QueryRow(
		`SELECT id, user_id, port_id, device_id, device_name, device_info, last_ping, status, connected_at, disconnected_at
		 FROM devices WHERE id = ?`, id,
	).Scan(&d.ID, &d.UserID, &d.PortID, &d.DeviceID, &d.DeviceName, &d.DeviceInfo, &d.LastPing, &d.Status, &connectedAt, &disconnectedAt)
	if err != nil {
		return nil, err
	}
	if connectedAt.Valid {
		d.ConnectedAt = &connectedAt.Time
	}
	if disconnectedAt.Valid {
		d.DisconnectedAt = &disconnectedAt.Time
	}
	return d, nil
}

func GetDeviceByPortAndDeviceID(portID int64, deviceID string) (*Device, error) {
	d := &Device{}
	var connectedAt, disconnectedAt sql.NullTime
	err := DB.QueryRow(
		`SELECT id, user_id, port_id, device_id, device_name, device_info, last_ping, status, connected_at, disconnected_at
		 FROM devices WHERE port_id = ? AND device_id = ?`, portID, deviceID,
	).Scan(&d.ID, &d.UserID, &d.PortID, &d.DeviceID, &d.DeviceName, &d.DeviceInfo, &d.LastPing, &d.Status, &connectedAt, &disconnectedAt)
	if err != nil {
		return nil, err
	}
	if connectedAt.Valid {
		d.ConnectedAt = &connectedAt.Time
	}
	if disconnectedAt.Valid {
		d.DisconnectedAt = &disconnectedAt.Time
	}
	return d, nil
}

func GetUserDevices(userID int64) ([]*Device, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, port_id, device_id, device_name, device_info, last_ping, status, connected_at, disconnected_at
		 FROM devices WHERE user_id = ? ORDER BY connected_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*Device
	for rows.Next() {
		d := &Device{}
		var connectedAt, disconnectedAt sql.NullTime
		err := rows.Scan(&d.ID, &d.UserID, &d.PortID, &d.DeviceID, &d.DeviceName, &d.DeviceInfo, &d.LastPing, &d.Status, &connectedAt, &disconnectedAt)
		if err != nil {
			return nil, err
		}
		if connectedAt.Valid {
			d.ConnectedAt = &connectedAt.Time
		}
		if disconnectedAt.Valid {
			d.DisconnectedAt = &disconnectedAt.Time
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func GetPortDevices(portID int64) ([]*Device, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, port_id, device_id, device_name, device_info, last_ping, status, connected_at, disconnected_at
		 FROM devices WHERE port_id = ? ORDER BY connected_at DESC`, portID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*Device
	for rows.Next() {
		d := &Device{}
		var connectedAt, disconnectedAt sql.NullTime
		err := rows.Scan(&d.ID, &d.UserID, &d.PortID, &d.DeviceID, &d.DeviceName, &d.DeviceInfo, &d.LastPing, &d.Status, &connectedAt, &disconnectedAt)
		if err != nil {
			return nil, err
		}
		if connectedAt.Valid {
			d.ConnectedAt = &connectedAt.Time
		}
		if disconnectedAt.Valid {
			d.DisconnectedAt = &disconnectedAt.Time
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func GetOnlineDevices(portID int64) ([]*Device, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, port_id, device_id, device_name, device_info, last_ping, status, connected_at, disconnected_at
		 FROM devices WHERE port_id = ? AND status = 'online' ORDER BY connected_at DESC`, portID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []*Device
	for rows.Next() {
		d := &Device{}
		var connectedAt, disconnectedAt sql.NullTime
		err := rows.Scan(&d.ID, &d.UserID, &d.PortID, &d.DeviceID, &d.DeviceName, &d.DeviceInfo, &d.LastPing, &d.Status, &connectedAt, &disconnectedAt)
		if err != nil {
			return nil, err
		}
		if connectedAt.Valid {
			d.ConnectedAt = &connectedAt.Time
		}
		if disconnectedAt.Valid {
			d.DisconnectedAt = &disconnectedAt.Time
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func (d *Device) SetOnline() error {
	_, err := DB.Exec(
		`UPDATE devices SET status = 'online', connected_at = CURRENT_TIMESTAMP WHERE id = ?`, d.ID,
	)
	return err
}

func (d *Device) SetOffline() error {
	_, err := DB.Exec(
		`UPDATE devices SET status = 'offline', disconnected_at = CURRENT_TIMESTAMP WHERE id = ?`, d.ID,
	)
	return err
}

func (d *Device) UpdatePing(ping int) error {
	d.LastPing = ping
	_, err := DB.Exec(`UPDATE devices SET last_ping = ? WHERE id = ?`, ping, d.ID)
	return err
}

func SetAllDevicesOffline() error {
	_, err := DB.Exec(`UPDATE devices SET status = 'offline', disconnected_at = CURRENT_TIMESTAMP WHERE status = 'online'`)
	return err
}

func SetPortDevicesOffline(portID int64) error {
	_, err := DB.Exec(`UPDATE devices SET status = 'offline', disconnected_at = CURRENT_TIMESTAMP WHERE port_id = ? AND status = 'online'`, portID)
	return err
}

// DeviceLog 设备日志
type DeviceLog struct {
	ID        int64     `json:"id"`
	DeviceID  int64     `json:"device_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func AddDeviceLog(deviceID int64, level, message string) error {
	_, err := DB.Exec(
		`INSERT INTO device_logs (device_id, level, message) VALUES (?, ?, ?)`,
		deviceID, level, message,
	)
	return err
}

func GetDeviceLogs(deviceID int64, limit int) ([]*DeviceLog, error) {
	rows, err := DB.Query(
		`SELECT id, device_id, level, message, created_at FROM device_logs
		 WHERE device_id = ? ORDER BY created_at DESC LIMIT ?`, deviceID, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*DeviceLog
	for rows.Next() {
		l := &DeviceLog{}
		err := rows.Scan(&l.ID, &l.DeviceID, &l.Level, &l.Message, &l.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}

	return logs, nil
}

func GetAllOnlineDeviceCount() (int, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM devices WHERE status = 'online'`).Scan(&count)
	return count, err
}

type DeviceWithPort struct {
	Device
	PortNumber int    `json:"port_number"`
	Username   string `json:"username"`
}

func GetAllDevicesWithPort(page, pageSize int, status string) ([]*DeviceWithPort, int, error) {
	var total int
	var countQuery string
	var queryArgs []interface{}

	if status != "" {
		countQuery = `SELECT COUNT(*) FROM devices WHERE status = ?`
		queryArgs = append(queryArgs, status)
	} else {
		countQuery = `SELECT COUNT(*) FROM devices`
	}

	err := DB.QueryRow(countQuery, queryArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	var query string
	if status != "" {
		query = `SELECT d.id, d.user_id, d.port_id, d.device_id, d.device_name, d.device_info, d.last_ping, d.status,
				 d.connected_at, d.disconnected_at, p.port, u.username
				 FROM devices d
				 LEFT JOIN ports p ON d.port_id = p.id
				 LEFT JOIN users u ON d.user_id = u.id
				 WHERE d.status = ?
				 ORDER BY d.connected_at DESC LIMIT ? OFFSET ?`
		queryArgs = append(queryArgs, pageSize, offset)
	} else {
		query = `SELECT d.id, d.user_id, d.port_id, d.device_id, d.device_name, d.device_info, d.last_ping, d.status,
				 d.connected_at, d.disconnected_at, p.port, u.username
				 FROM devices d
				 LEFT JOIN ports p ON d.port_id = p.id
				 LEFT JOIN users u ON d.user_id = u.id
				 ORDER BY d.connected_at DESC LIMIT ? OFFSET ?`
		queryArgs = []interface{}{pageSize, offset}
	}

	rows, err := DB.Query(query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var devices []*DeviceWithPort
	for rows.Next() {
		d := &DeviceWithPort{}
		var connectedAt, disconnectedAt sql.NullTime
		var portNumber sql.NullInt64
		var username sql.NullString
		err := rows.Scan(&d.ID, &d.UserID, &d.PortID, &d.DeviceID, &d.DeviceName, &d.DeviceInfo, &d.LastPing, &d.Status,
			&connectedAt, &disconnectedAt, &portNumber, &username)
		if err != nil {
			return nil, 0, err
		}
		if connectedAt.Valid {
			d.ConnectedAt = &connectedAt.Time
		}
		if disconnectedAt.Valid {
			d.DisconnectedAt = &disconnectedAt.Time
		}
		if portNumber.Valid {
			d.PortNumber = int(portNumber.Int64)
		}
		if username.Valid {
			d.Username = username.String
		}
		devices = append(devices, d)
	}

	return devices, total, nil
}
