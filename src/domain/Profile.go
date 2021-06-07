package domain

type Profile struct {
	Id string
	Username string `json:"username"`
	ProfilePhoto string `json:"profilePhoto"`
}
