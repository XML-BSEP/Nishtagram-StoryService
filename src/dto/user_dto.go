package dto

type UserDTO struct {
	UserId string `json:"userid" validate:"required"`
}
