package main

import (
	"strconv"

	"github.com/theplant/luhn"
)

func ValidateUserInfo(input *UserInfo) error {
	if input.Password == "" || input.Username == "" {
		return ErrValidation
	}
	return nil
}

func ValidateLogin(input *UserInfo, existingInfo *UserInfo) error {
	if input.Username == existingInfo.Username {
		if input.Password == existingInfo.Password {
			return nil
		}
	}
	return ErrValidation
}

func ValidateOrder(cursor *Cursor, username string, orderToValidate string) error {
	userOrders, _ := cursor.GetOrders()
	for _, order := range userOrders {
		if orderToValidate == order.Number {
			orderInt, err := strconv.Atoi(orderToValidate)
			if err != nil {
				return err
			}
			if luhn.Valid(orderInt) {
				return nil
			} else {
				return ErrValidation
			}

		}
	}
	return ErrValidation
}
