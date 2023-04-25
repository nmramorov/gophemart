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
