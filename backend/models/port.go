package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Port struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Port      int       `json:"port"`
	Protocol  string    `json:"protocol"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	Status    string    `json:"status"`
	ExpireAt  time.Time `json:"expire_at"`
	CreatedAt time.Time `json:"created_at"`
}

func generateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func CreatePort(userID int64, port int, protocol string, days int) (*Port, error) {
	token := generateToken()
	expireAt := time.Now().AddDate(0, 0, days)

	result, err := DB.Exec(
		`INSERT INTO ports (user_id, port, protocol, token, status, expire_at) VALUES (?, ?, ?, ?, 'active', ?)`,
		userID, port, protocol, token, expireAt,
	)
	if err != nil {
		return nil, err
	}

	id, _ := result.LastInsertId()
	return GetPortByID(id)
}

func GetPortByID(id int64) (*Port, error) {
	p := &Port{}
	err := DB.QueryRow(
		`SELECT id, user_id, port, protocol, name, token, status, expire_at, created_at FROM ports WHERE id = ?`, id,
	).Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPortByToken(token string) (*Port, error) {
	p := &Port{}
	err := DB.QueryRow(
		`SELECT id, user_id, port, protocol, name, token, status, expire_at, created_at FROM ports WHERE token = ?`, token,
	).Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetPortByNumber(port int) (*Port, error) {
	p := &Port{}
	err := DB.QueryRow(
		`SELECT id, user_id, port, protocol, name, token, status, expire_at, created_at FROM ports WHERE port = ?`, port,
	).Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func GetUserPorts(userID int64) ([]*Port, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, port, protocol, name, token, status, expire_at, created_at
		 FROM ports WHERE user_id = ? ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []*Port
	for rows.Next() {
		p := &Port{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		ports = append(ports, p)
	}

	return ports, nil
}

func GetActivePorts() ([]*Port, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, port, protocol, name, token, status, expire_at, created_at
		 FROM ports WHERE status = 'active'`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []*Port
	for rows.Next() {
		p := &Port{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		ports = append(ports, p)
	}

	return ports, nil
}

func (p *Port) Update() error {
	_, err := DB.Exec(
		`UPDATE ports SET protocol = ?, name = ?, status = ?, expire_at = ? WHERE id = ?`,
		p.Protocol, p.Name, p.Status, p.ExpireAt, p.ID,
	)
	return err
}

func (p *Port) Delete() error {
	// 先删除关联的设备
	_, err := DB.Exec(`DELETE FROM devices WHERE port_id = ?`, p.ID)
	if err != nil {
		return err
	}
	_, err = DB.Exec(`DELETE FROM ports WHERE id = ?`, p.ID)
	return err
}

func (p *Port) Renew(days int) error {
	// 如果已过期，从当前时间开始计算；否则从过期时间开始
	baseTime := p.ExpireAt
	if baseTime.Before(time.Now()) {
		baseTime = time.Now()
	}
	p.ExpireAt = baseTime.AddDate(0, 0, days)
	p.Status = "active"
	return p.Update()
}

func (p *Port) RegenerateToken() error {
	p.Token = generateToken()
	_, err := DB.Exec(`UPDATE ports SET token = ? WHERE id = ?`, p.Token, p.ID)
	return err
}

func IsPortUsed(port int) (bool, error) {
	var count int
	err := DB.QueryRow(`SELECT COUNT(*) FROM ports WHERE port = ?`, port).Scan(&count)
	return count > 0, err
}

func GetExpiredPorts() ([]*Port, error) {
	rows, err := DB.Query(
		`SELECT id, user_id, port, protocol, name, token, status, expire_at, created_at
		 FROM ports WHERE status = 'active' AND expire_at < CURRENT_TIMESTAMP`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ports []*Port
	for rows.Next() {
		p := &Port{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		ports = append(ports, p)
	}

	return ports, nil
}

func GetAllPorts(page, pageSize int) ([]*Port, int, error) {
	var total int
	err := DB.QueryRow(`SELECT COUNT(*) FROM ports`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		`SELECT id, user_id, port, protocol, name, token, status, expire_at, created_at
		 FROM ports ORDER BY id DESC LIMIT ? OFFSET ?`, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ports []*Port
	for rows.Next() {
		p := &Port{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt)
		if err != nil {
			return nil, 0, err
		}
		ports = append(ports, p)
	}

	return ports, total, nil
}

type PortWithUser struct {
	Port
	Username string `json:"username"`
}

func GetAllPortsWithUser(page, pageSize int) ([]*PortWithUser, int, error) {
	var total int
	err := DB.QueryRow(`SELECT COUNT(*) FROM ports`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	rows, err := DB.Query(
		`SELECT p.id, p.user_id, p.port, p.protocol, p.name, p.token, p.status, p.expire_at, p.created_at, u.username
		 FROM ports p LEFT JOIN users u ON p.user_id = u.id ORDER BY p.id DESC LIMIT ? OFFSET ?`, pageSize, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var ports []*PortWithUser
	for rows.Next() {
		p := &PortWithUser{}
		err := rows.Scan(&p.ID, &p.UserID, &p.Port, &p.Protocol, &p.Name, &p.Token, &p.Status, &p.ExpireAt, &p.CreatedAt, &p.Username)
		if err != nil {
			return nil, 0, err
		}
		ports = append(ports, p)
	}

	return ports, total, nil
}
