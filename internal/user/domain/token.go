package domain

type Token struct {
	AccessToken      string
	expiresIn        int
	RefreshToken     string
	RefreshExpiresIn int
}
