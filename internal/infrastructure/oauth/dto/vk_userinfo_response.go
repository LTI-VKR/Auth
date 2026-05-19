package dto

type VkUserinfoResponse struct {
	User struct {
		UserId    string `json:"user_id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
		Email     string `json:"email"`
		Sex       int    `json:"sex"`
		Verified  bool   `json:"verified"`
		Birthday  string `json:"birthday"`
	} `json:"user"`
}
