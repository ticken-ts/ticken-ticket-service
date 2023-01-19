package security

type JWTVerifier interface {
	Verify(rawJWT string) (*JWT, error)
}
