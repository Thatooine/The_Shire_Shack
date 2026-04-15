package authentication

type LoginClaim struct {
	UserID         string `json:"userID"`
	ExpirationTime int64  `json:"expirationTime"`
	LastName       string `json:"lastName"`
	Email          string `json:"email"`
}
