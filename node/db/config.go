package db

import (
	"fmt"
)

// ConfigType config type
type ConfigType string

const (
	// ConfigTronHeight is used for storing the height of scanned blocks.
	ConfigTronHeight ConfigType = "tron_height"
)

// SaveConfigValue saves a configuration value.
func (d *SQLDB) SaveConfigValue(key ConfigType, value string) error {
	query := fmt.Sprintf(
		`INSERT INTO %s (name, value) VALUES (?, ?)
				ON DUPLICATE KEY UPDATE value=?`, configTable)
	_, err := d.db.Exec(query, key, value, value)

	return err
}

// LoadConfigValue loads a configuration value.
func (d *SQLDB) LoadConfigValue(key ConfigType, out interface{}) error {
	query := fmt.Sprintf("SELECT value FROM %s WHERE name=?", configTable)
	return d.db.Get(out, query, key)
}
