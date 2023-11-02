package model

type UpdatePassword struct {
	Tel         string `json:"tel"`
	Otp         string `json:"otp"`
	Password    string `json:"current_password"`
	NewPassword string `json:"new_password"`
}

type UserResponse struct {
	UUID           string `json:"uuid"`
	HashedPassword string `json:"hashed_password"`
	IsPinSet       bool   `json:"is_pin_set"`
	IsProfileSet   bool   `json:"is_profile_set"`
}
