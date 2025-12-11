package models

type Log struct {
	ID        uint   `json:"id" autoIncrement:"true"`
	Username  string `json:"username"`
	UserId    uint   `json:"user_id"`
	RoleID    uint   `json:"role_id"`
	LoginDate string `json:"login_date"`
}
