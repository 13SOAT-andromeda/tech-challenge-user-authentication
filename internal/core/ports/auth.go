package ports

type LoginInput struct {
	Document string `json:"document"`
	Password string `json:"password"`
}

type UserOutput struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Contact string `json:"contact"`
	Role    string `json:"role"`
}

type LoginOutput struct {
	AccessToken  string     `json:"access_token"`
	RefreshToken string     `json:"refresh_token"`
	ExpiresIn    int64      `json:"expires_in"`
	JTI          string     `json:"jti"`
	User         UserOutput `json:"user"`
}
