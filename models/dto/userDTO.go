package dto

// struct parametro para os endpoinds de post
type UserRequest struct {
	Email     string `json:"email"`
	SenhaHash string `json:"senhaHash"`
}
