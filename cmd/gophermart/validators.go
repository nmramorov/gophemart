package main


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
			return nil
		}
	}
	return ErrValidation
}
