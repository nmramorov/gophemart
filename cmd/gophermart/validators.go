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

func ValidateOrder(cursor *Cursor, newOrder *Order) error {
	orders := cursor.GetAllOrders()
	for _, order := range orders {
		if order.Number == newOrder.Number {
			return ErrValidation
		}
	}
	return nil
}
