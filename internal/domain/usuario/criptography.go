package usuario

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hashedPassword, password string) bool
}

type TokenGenerator interface {
	GenerateToken(userId string, role string) (string, error)
}
