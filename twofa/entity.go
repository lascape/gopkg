package twofa

type Authenticator struct {
	Secret string
	Expire int
	Code   string
	Error  error
}
