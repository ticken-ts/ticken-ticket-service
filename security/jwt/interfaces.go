package jwt

type Verifier interface {
	Verify(rawJWT string) (*Token, error)
}
