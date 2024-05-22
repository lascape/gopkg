package cryptox

type cipherError string

func (e cipherError) Error() string {
	return string(e)
}
