package hasher

type Hasher interface {
	Hash(plaintext string) (hash string, err error)
	Match(plaintext, hash string) bool
}
