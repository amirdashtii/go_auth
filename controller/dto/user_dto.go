package dto

type UserProfileResponse struct {
    ID        string `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
   
}

type UserUpdateRequest struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
} 

type ChangePasswordRequest struct {
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required,min=8"`
}
