package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ahmadmirdas/julo-test/repository/database/models"
	"github.com/ahmadmirdas/julo-test/server/middleware"
	"github.com/ahmadmirdas/julo-test/utils/activity"
	"github.com/ahmadmirdas/julo-test/utils/log"
	"github.com/ahmadmirdas/julo-test/utils/response"
	"github.com/golang-jwt/jwt/v4"
)

const (
	Limit = 10
)

type handlerWallet struct {
	walletRepo models.WalletDBRepo
}

type HandlerWallet interface {
	InitAccountWallet(w http.ResponseWriter, r *http.Request)
	EnableWallet(w http.ResponseWriter, r *http.Request)
	ViewWalletBalance(w http.ResponseWriter, r *http.Request)
	DepositWallet(w http.ResponseWriter, r *http.Request)
	WithdrawWallet(w http.ResponseWriter, r *http.Request)
	DisableWallet(w http.ResponseWriter, r *http.Request)
}

func NewHandlerWallet(walletRepo models.WalletDBRepo) HandlerWallet {
	return &handlerWallet{
		walletRepo: walletRepo,
	}
}

func (h *handlerWallet) InitAccountWallet(w http.ResponseWriter, r *http.Request) {
	ctx := activity.NewContext("Handler.InitAccountWallet")
	customerXId := r.FormValue("customer_xid")
	if customerXId == "" {
		log.WithContext(ctx).Warn("[Handler InitAccountWallet] customer xid is required")
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusBadRequest,
				Message: "customer_xid is required",
			},
		}, http.StatusBadRequest)
		return
	}
	token, err := middleware.GenerateToken(customerXId)
	if err != nil {
		log.WithContext(ctx).Warn("[Handler InitAccountWallet] Error when generate token")
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	httpResponseWrite(w, response.ResponseAPI{
		Status: "success",
		Data: map[string]string{
			"token": token,
		},
	}, http.StatusOK)
}

func (h *handlerWallet) EnableWallet(w http.ResponseWriter, r *http.Request) {
	ctx := activity.NewContext("Handler.EnableWallet")
	cus := r.Context().Value(middleware.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	res, err := h.walletRepo.EnableWallet(custXId)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler EnableWallet] error when enable wallet, error: ", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	status := "enabled"
	if !res.IsEnabled {
		status = "disabled"
	}
	httpResponseWrite(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWallet{
			ID:        res.ID,
			OwnedBy:   res.OwnedBy,
			Status:    status,
			EnabledAt: res.EnabledAt.String(),
			Balance:   res.Balance,
		},
	}, http.StatusOK)
}

func (h *handlerWallet) ViewWalletBalance(w http.ResponseWriter, r *http.Request) {
	ctx := activity.NewContext("Handler.ViewWalletBalance")
	cus := r.Context().Value(middleware.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	wallet, err := h.walletRepo.GetWallet(custXId)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler ViewWalletBalance] error when query get wallet, error: ", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	status := "enabled"
	if !wallet.IsEnabled {
		log.WithContext(ctx).Error("[Handler ViewWalletBalance] your wallet is disabled, cannot view")
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: "your wallet is disabled, cannot view",
			},
		}, http.StatusInternalServerError)
		return
	}
	httpResponseWrite(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWallet{
			ID:        wallet.ID,
			OwnedBy:   wallet.OwnedBy,
			Status:    status,
			EnabledAt: wallet.EnabledAt.String(),
			Balance:   wallet.Balance,
		},
	}, http.StatusOK)
}

func (h *handlerWallet) DepositWallet(w http.ResponseWriter, r *http.Request) {
	ctx := activity.NewContext("Handler.DepositWallet")
	cus := r.Context().Value(middleware.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	amountStr := r.FormValue("amount")
	referenceId := r.FormValue("reference_id")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.WithContext(ctx).Error("[Handler DepositWallet] invalid amount, error : %v", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	wallet, err := h.walletRepo.GetWallet(custXId)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler DepositWallet] error when query get wallet, error: ", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	if !wallet.IsEnabled {
		log.WithContext(ctx).Error("[Handler DepositWallet] your wallet is disabled, cannot deposit")
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: "your wallet is disabled, cannot deposit",
			},
		}, http.StatusInternalServerError)
		return
	}

	param := models.ParamWalletDeposit{
		WalletID:    wallet.ID,
		Balance:     wallet.Balance,
		CustomerXId: custXId,
		Amount:      amount,
		ReferenceID: referenceId,
	}
	res, err := h.walletRepo.WalletDeposit(param)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler DepositWallet] error when query wallet deposit, error: %v", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	httpResponseWrite(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseDepositWallet{
			ID:          res.ID,
			DepositedBy: res.Wallet.OwnedBy,
			Status:      res.Status,
			DepositAt:   res.CreatedAt.String(),
			Amount:      res.Amount,
			ReferenceId: res.ReferenceID,
		},
	}, http.StatusOK)
}

func (h *handlerWallet) WithdrawWallet(w http.ResponseWriter, r *http.Request) {
	ctx := activity.NewContext("Handler.WithdrawWallet")
	cus := r.Context().Value(middleware.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	amountStr := r.FormValue("amount")
	referenceId := r.FormValue("reference_id")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.WithContext(ctx).Error("[Handler WithdrawWallet] invalid amount, error : %v", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	wallet, err := h.walletRepo.GetWallet(custXId)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler WithdrawWallet] error when query get wallet, error: ", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	if !wallet.IsEnabled {
		log.WithContext(ctx).Error("[Handler WithdrawWallet] your wallet is disabled, cannot withdraw")
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: "your wallet is disabled, cannot withdraw",
			},
		}, http.StatusInternalServerError)
		return
	}

	param := models.ParamWalletWithdraw{
		WalletID:    wallet.ID,
		Balance:     wallet.Balance,
		CustomerXId: custXId,
		Amount:      amount,
		ReferenceID: referenceId,
	}
	res, err := h.walletRepo.WalletWithdraw(param)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler WithdrawWallet] error when query withdraw wallet, error: %v", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	httpResponseWrite(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWithdrawWallet{
			ID:          res.ID,
			WithdrawnBy: res.Wallet.OwnedBy,
			Status:      res.Status,
			WithdrawnAt: res.CreatedAt.String(),
			Amount:      res.Amount,
			ReferenceId: res.ReferenceID,
		},
	}, http.StatusOK)
}

func (h *handlerWallet) DisableWallet(w http.ResponseWriter, r *http.Request) {
	ctx := activity.NewContext("Handler.DisableWallet")
	cus := r.Context().Value(middleware.Customer).(jwt.MapClaims)
	custXId := cus["customer_xid"].(string)

	formDisabled := r.FormValue("is_disabled")
	isDisabled, err := strconv.ParseBool(formDisabled)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler DisableWallet] error parse is_disabled to bool, error: ", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	var isEnabled bool
	if !isDisabled {
		isDisabled = true
	}
	res, err := h.walletRepo.UpdateStatusWallet(custXId, isEnabled)
	if err != nil {
		log.WithContext(ctx).Errorf("[Handler DisableWallet] error when enable wallet, error: ", err)
		httpResponseWrite(w, response.ResponseAPI{
			Error_: &response.ApiError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	status := "enabled"
	if !res.IsEnabled {
		status = "disabled"
	}
	httpResponseWrite(w, response.ResponseAPI{
		Status: "success",
		Data: ResponseWallet{
			ID:        res.ID,
			OwnedBy:   res.OwnedBy,
			Status:    status,
			EnabledAt: res.DisabledAt.String(),
			Balance:   res.Balance,
		},
	}, http.StatusOK)
}

func httpResponseWrite(rw http.ResponseWriter, data interface{}, statusCode int) {
	rw.Header().Set("Content-type", "application/json")
	rw.WriteHeader(statusCode)
	json.NewEncoder(rw).Encode(data)
}
