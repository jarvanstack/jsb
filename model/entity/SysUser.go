package entity
type SysUser struct {
	UserId int64 `json:"user_id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Nickname string `json:"nickname"`
	Avatar string `json:"avatar"`
}