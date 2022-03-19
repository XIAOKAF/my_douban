package model

type OAuth struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type Token struct {
	AccessToken string
}
