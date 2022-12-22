package entity

import "time"

type Wallet struct {
	tableName  struct{}  `pg:"wallet"`
	ID         string    `json:"id" pg:"id,pk"`
	OwnedBy    string    `json:"-"  pg:"owned_by"`
	IsEnabled  bool      `json:"-"  pg:"is_enabled"`
	Balance    float64   `json:"-"  pg:"balance"`
	EnabledAt  time.Time `json:"-"  pg:"enabled_at"`
	DisabledAt time.Time `json:"-"  pg:"disabled_at"`
}
