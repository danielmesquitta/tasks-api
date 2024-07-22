package symcrypt

type SymmetricalEncrypter interface {
	Encrypt(plaintext string) (encrypted string, err error)
	Decrypt(encrypted string) (plaintext string, err error)
}
