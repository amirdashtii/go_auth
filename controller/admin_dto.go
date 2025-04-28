package controller

type AdminUserListResponse struct {
	Users []AdminUserResponse `json:"users"`
}

type AdminUserResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Role      string `json:"role"`
}

type AdminUserUpdateRequest struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	Status    *string `json:"status,omitempty"`
	Role      *string `json:"role,omitempty"`
}
