package api

import (
	"github.com/nmramorov/gophemart/internal/db"
	"github.com/nmramorov/gophemart/internal/errors"
	"github.com/nmramorov/gophemart/internal/models"
)


func ValidateUserInfo(input *models.UserInfo) error {
	if input.Password == "" || input.Username == "" {
		return errors.ErrValidation
	}
	return nil
}

func ValidateLogin(input *models.UserInfo, existingInfo *models.UserInfo) error {
	if input.Username == existingInfo.Username {
		if input.Password == existingInfo.Password {
			return nil
		}
	}
	return errors.ErrValidation
}

func ValidateOrder(cursor *db.Cursor, newOrder *models.Order) error {
	orders := cursor.GetAllOrders()
	for _, order := range orders {
		if order.Number == newOrder.Number {
			return errors.ErrValidation
		}
	}
	return nil
}
