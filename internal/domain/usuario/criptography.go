package usuario

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
}
