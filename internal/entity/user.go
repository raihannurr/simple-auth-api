package entity

type User struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	Email       string `json:"email"`
	Verified    bool   `json:"verified"`
	Description string `json:"description"`
}
