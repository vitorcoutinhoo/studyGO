package apierr

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestBindingMessageIdentificaCampoERegra(t *testing.T) {
	type request struct {
		Email string `json:"email" validate:"required,email"`
		Senha string `json:"senha" validate:"min=6"`
	}
	validate := validator.New()

	for _, teste := range []struct {
		valor    request
		esperado string
	}{
		{request{}, `O campo “email” é obrigatório.`},
		{request{Email: "inválido"}, `O campo “email” deve conter um e-mail válido.`},
		{request{Email: "usuario@example.com", Senha: "123"}, `O campo “senha” deve ter no mínimo 6 caracteres.`},
	} {
		err := validate.Struct(teste.valor)
		if mensagem := bindingMessage(err); mensagem != teste.esperado {
			t.Fatalf("mensagem = %q, esperado %q", mensagem, teste.esperado)
		}
	}
}

func TestBindingMessageNaoExpoeErroTecnico(t *testing.T) {
	if mensagem := bindingMessage(errors.New("erro interno do decoder")); mensagem != "Não foi possível interpretar os dados enviados. Verifique os campos e tente novamente." {
		t.Fatalf("mensagem = %q", mensagem)
	}
}
