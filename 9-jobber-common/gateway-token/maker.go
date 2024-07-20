package token

type TokenMaker interface {
	CreateToken(id string) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}
