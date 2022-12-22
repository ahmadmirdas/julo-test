package entity

import "time"

type History struct {
	tableName   struct{}  `pg:"history"`
	ID          string    `json:"id" pg:"id,pk"`
	WalletID    string    `json:"-"  pg:"wallet_id"`
	Wallet      *Wallet   `json:"-"  pg:"fk:wallet_id"`
	Status      string    `json:"-"  pg:"status"`
	Amount      float64   `json:"-"  pg:"amount"`
	Type        string    `json:"-"  pg:"type"`
	ReferenceID string    `json:"-"  pg:"reference_id"`
	CreatedAt   time.Time `json:"-"  pg:"created_at"`
}
