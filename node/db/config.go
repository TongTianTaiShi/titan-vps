package db

import (
	"fmt"
)

// ConfigType config type
type ConfigType string

const (
	// ConfigTronHeight
	ConfigTronHeight ConfigType = "tron_height"
)

// SaveConfigValue save config value
func (n *SQLDB) SaveConfigValue(key ConfigType, value string) error {
	// update record table
	query := fmt.Sprintf(
		`INSERT INTO %s (name, value) VALUES (?, ?)
				ON DUPLICATE KEY UPDATE value=?`, configTable)
	_, err := n.db.Exec(query, key, value, value)

	return err
}

// LoadConfigValue load config value
func (n *SQLDB) LoadConfigValue(key ConfigType, out interface{}) error {
	query := fmt.Sprintf("SELECT value FROM %s WHERE name=?", configTable)
	return n.db.Get(out, query, key)
}
