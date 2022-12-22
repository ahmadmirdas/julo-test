package models

var (
	WalletStatusEnabled  string = "enabled"
	WalletStatusDisabled string = "disabled"
)

type ParamWalletDeposit struct {
	WalletID    string
	Balance     float64
	Amount      float64
	CustomerXId string
	ReferenceID string
}

type ParamWalletWithdraw struct {
	WalletID    string
	Balance     float64
	Amount      float64
	CustomerXId string
	ReferenceID string
}
