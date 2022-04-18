package user

type SendCode struct {
	Phone        string `json:"phone" binding:"required"`
	CaptchaId    string `json:"captchaId" binding:"required"`
	CaptchaValue string `json:"captchaValue" binding:"required"`
}

type CheckPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

type CheckPhoneResponse struct {
	Phone string `json:"phone"`
	Exist bool   `json:"exist"`
}

type RegisterRequest struct {
	Phone    string `json:"phone" binding:"required"`
	Password string `json:"password" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type LoginGeneralRequest struct {
	Info     string `json:"info" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	Phone         string `json:"phone" binding:"required"`
	Password      string `json:"password" binding:"required"`
	PasswordAgain string `json:"passwordAgain" binding:"required"`
	Code          string `json:"code" binding:"required"`
}

type UpdateNickNameRequest struct {
	NickName string `json:"nickName" binding:"required"`
}

type UpdatePhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type UpdateAvatarRequest struct {
	AvatarUrl string `json:"avatarUrl" binding:"required"`
}

type UpdateSexRequest struct {
	Sex int `json:"sex" binding:"required"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type UpdateSignatureRequest struct {
	Signature string `json:"signature" binding:"required"`
}

type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required"`
}
