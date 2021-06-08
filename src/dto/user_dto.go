package dto

type UserDTO struct {
	UserId string `json:"id" validate:"required"`
}
