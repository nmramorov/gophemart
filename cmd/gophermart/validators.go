package main


func ValidateUserRegistrationInfo(input *UserRegistrationInfo) error {
	if input.Password == "" || input.Username == "" {
		return ErrValidation
	}
	return nil
}