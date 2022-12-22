package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ahmadmirdas/julo-test/repository/database/models/entity"
	"github.com/go-pg/pg/v10"
)

type WalletDBRepo interface {
	EnableWallet(customerXId string) (*entity.Wallet, error)
	GetWallet(customerXId string) (*entity.Wallet, error)
	WalletDeposit(param ParamWalletDeposit) (*entity.History, error)
	WalletWithdraw(param ParamWalletWithdraw) (*entity.History, error)
	UpdateStatusWallet(customerXId string, status bool) (*entity.Wallet, error)
}

type dbWalletRepo struct {
	mutex  sync.RWMutex
	dbConn *pg.DB
}

func NewDBWalletRepo(c *pg.DB) WalletDBRepo {
	return &dbWalletRepo{dbConn: c}
}

func (p *dbWalletRepo) EnableWallet(customerXId string) (*entity.Wallet, error) {
	if customerXId == "" {
		return nil, errors.New("customerXId is empty")
	}

	resWallet, err := p.GetWallet(customerXId)
	if err != nil {
		return nil, err
	}

	if resWallet != nil {
		if resWallet.IsEnabled {
			return nil, fmt.Errorf("wallet is already enabled")
		}
	}

	wallet := entity.Wallet{
		OwnedBy:   customerXId,
		IsEnabled: true,
		EnabledAt: time.Now(),
	}
	p.mutex.Lock()
	res, err := p.dbConn.Model(&wallet).
		OnConflict("(owned_by) DO UPDATE").
		Set("is_enabled = EXCLUDED.is_enabled").
		Insert()
	if err != nil {
		return nil, err
	}
	p.mutex.Unlock()
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("failed enabled wallet")
	}

	return &wallet, nil
}

func (p *dbWalletRepo) GetWallet(customerXId string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	p.mutex.Lock()
	err := p.dbConn.Model(&wallet).
		Where("owned_by = ?", customerXId).
		Select()
	if err != nil {
		if err != pg.ErrNoRows {
			return nil, err
		}
	}
	p.mutex.Unlock()
	return &wallet, nil
}

func (p *dbWalletRepo) WalletDeposit(param ParamWalletDeposit) (*entity.History, error) {
	var result entity.History

	wallet := entity.Wallet{}
	res, err := p.dbConn.Model(&wallet).
		Where("owned_by = ? ", param.CustomerXId).
		Set("balance = ?", param.Balance+param.Amount).
		Update()
	if err != nil {
		return nil, err
	}

	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("deposit failed - error update balance")
	}

	history := entity.History{
		WalletID:    param.WalletID,
		Status:      "success",
		Amount:      param.Amount,
		Type:        "deposit",
		ReferenceID: param.ReferenceID,
	}
	_, err = p.dbConn.Model(&history).Insert()
	if err != nil {
		return nil, err
	}

	err = p.dbConn.Model(&result).Relation("Wallet").
		Where("history.id = ?", history.ID).
		Select()
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (p *dbWalletRepo) WalletWithdraw(param ParamWalletWithdraw) (*entity.History, error) {
	var result entity.History

	wallet := entity.Wallet{}
	p.mutex.Lock()
	res, err := p.dbConn.Model(&wallet).
		Where("owned_by = ? ", param.CustomerXId).
		Set("balance = ?", param.Balance-param.Amount).
		Update()
	if err != nil {
		return nil, err
	}
	p.mutex.Unlock()

	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("withdraw failed - error update balance")
	}

	history := entity.History{
		WalletID:    param.WalletID,
		Status:      "success",
		Type:        "withdraw",
		Amount:      param.Amount,
		ReferenceID: param.ReferenceID,
	}
	p.mutex.Lock()
	_, err = p.dbConn.Model(&history).Insert()
	if err != nil {
		return nil, err
	}
	p.mutex.Unlock()

	p.mutex.RLock()
	err = p.dbConn.Model(&result).Relation("Wallet").
		Where("history.id = ?", history.ID).
		Select()
	if err != nil {
		return nil, err
	}
	p.mutex.RUnlock()

	return &result, nil
}

func (p *dbWalletRepo) UpdateStatusWallet(customerXId string, status bool) (*entity.Wallet, error) {
	wallet := entity.Wallet{}
	p.mutex.Lock()
	res, err := p.dbConn.Model(&wallet).
		Where("owned_by = ?", customerXId).
		Set("is_enabled = ?", status).
		Set("disabled_at = ?", time.Now()).
		Update()
	if err != nil {
		return nil, err
	}
	p.mutex.Unlock()
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("failed update status wallet")
	}

	resWallet, err := p.GetWallet(customerXId)
	if err != nil {
		return nil, err
	}

	return resWallet, nil
}
