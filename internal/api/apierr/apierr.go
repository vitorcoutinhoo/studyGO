package apierr

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"plantao/internal/domain/cargo"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/convite"
	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/setor"
	"plantao/internal/domain/shared"
	"plantao/internal/domain/smtp"
	"plantao/internal/domain/usuario"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Respond(ctx *gin.Context, err error) {
	status, resp := classify(err)
	ctx.JSON(status, resp)
}

// RespondBinding traduz erros de corpo e validação para mensagens seguras e
// úteis ao consumidor da API, sem expor detalhes internos do framework.
func RespondBinding(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, ErrorResponse{Code: "BAD_REQUEST", Message: bindingMessage(err)})
}

func bindingMessage(err error) string {
	if errors.Is(err, io.EOF) {
		return "O corpo da requisição é obrigatório."
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) && len(validationErrors) > 0 {
		fieldError := validationErrors[0]
		field := jsonFieldName(fieldError)
		switch fieldError.Tag() {
		case "required":
			return fmt.Sprintf("O campo “%s” é obrigatório.", field)
		case "email":
			return fmt.Sprintf("O campo “%s” deve conter um e-mail válido.", field)
		case "uuid":
			return fmt.Sprintf("O campo “%s” deve conter um UUID válido.", field)
		case "min":
			return fmt.Sprintf("O campo “%s” deve ter no mínimo %s caracteres.", field, fieldError.Param())
		case "max":
			return fmt.Sprintf("O campo “%s” deve ter no máximo %s caracteres.", field, fieldError.Param())
		case "oneof":
			return fmt.Sprintf("O campo “%s” possui um valor inválido.", field)
		default:
			return fmt.Sprintf("O campo “%s” é inválido.", field)
		}
	}

	var typeError *json.UnmarshalTypeError
	if errors.As(err, &typeError) {
		return fmt.Sprintf("O campo “%s” possui formato inválido.", typeError.Field)
	}
	var syntaxError *json.SyntaxError
	if errors.As(err, &syntaxError) {
		return "O corpo da requisição contém JSON inválido."
	}
	return "Não foi possível interpretar os dados enviados. Verifique os campos e tente novamente."
}

func jsonFieldName(fieldError validator.FieldError) string {
	var result strings.Builder
	for index, char := range fieldError.Field() {
		if index > 0 && char >= 'A' && char <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

func classify(err error) (int, ErrorResponse) {
	if isAny(err, financeiro.ErrorAtualizacaoVazia, plantao.ErrorAtualizacaoPlantaoVazia, smtp.ErrorAtualizacaoSMTPVazia) {
		return http.StatusBadRequest, ErrorResponse{Code: "BAD_REQUEST", Message: err.Error()}
	}

	// 404 — não encontrado
	if isAny(err,
		cargo.ErrorCargoNotFound,
		colaborador.ErrorColaboradorNotFound,
		comunicacao.ErrorModeloComunicacaoNotFound,
		convite.ErrorConviteNotFound,
		financeiro.ErrorFeriadoNotFound,
		financeiro.ErrorValorDiaNotFound,
		plantao.ErrorPlantaoNotFinded,
		plantao.ErrorColaboradorNotFound,
		plantao.ErrorPagamentoNotFound,
		setor.ErrorSetorNotFound,
		smtp.ErrorConfiguracaoSMTPNotFound,
		usuario.ErrorUserNotFound,
	) {
		return http.StatusNotFound, ErrorResponse{Code: "NOT_FOUND", Message: err.Error()}
	}

	// 409 — conflito
	if isAny(err,
		cargo.ErrorCargoAlreadyExists,
		colaborador.ErrorEmailAlreadyExists,
		colaborador.ErrorTelefoneAlreadyExists,
		colaborador.ErrorInvalidEmail,
		comunicacao.ErrorModeloComunicacaoAlreadyExists,
		setor.ErrorSetorAlreadyExists,
		usuario.ErrorEmailAlreadyExists,
		usuario.ErrorEmailexists,
		plantao.ErrorExistingPlantao,
		plantao.ErrorPlantaoJaFechado,
		plantao.ErrorPlantaoJaPago,
		plantao.ErrorDetalhesExistentes,
		plantao.ErrorPagamentoExistente,
		plantao.ErrorPagamentoInconsistente,
		plantao.ErrorDetalhesInconsistentes,
		plantao.ErrorConflitoConcorrencia,
		plantao.ErrorPlantaoNaoEditavel,
		financeiro.ErrorValorDiaAlreadyExists,
		financeiro.ErrorConflitoValorDia,
		smtp.ErrorConfiguracaoSMTPAlreadyExists,
	) {
		return http.StatusConflict, ErrorResponse{Code: "CONFLICT", Message: err.Error()}
	}

	// 410 — expirado/consumido
	if isAny(err,
		convite.ErrorConviteExpired,
		convite.ErrorConviteUsed,
	) {
		return http.StatusGone, ErrorResponse{Code: "GONE", Message: err.Error()}
	}

	// 422 — regra de negócio / validação de domínio
	if isAny(err,
		cargo.ErrorCargoNomeInvalido,
		colaborador.ErrorInvalidTelefone,
		colaborador.ErrorInvalidStatus,
		colaborador.ErrorInvalidCargo,
		colaborador.ErrorInvalidSetor,
		colaborador.ErrorInactiveColaborador,
		comunicacao.ErrorInvalidNome,
		comunicacao.ErrorInvalidTipoComunicacao,
		comunicacao.ErrorAssunto,
		comunicacao.ErrorInvalidCorpo,
		comunicacao.ErrorInvalidStatus,
		financeiro.ErrorFeriadoDataInvalid,
		financeiro.ErrorFeriadoNotMunicipal,
		financeiro.ErrorTipoDiaInvalido,
		financeiro.ErrorValorDiaInvalido,
		financeiro.ErrorPrecisaoValorDia,
		financeiro.ErrorLimiteValorDia,
		plantao.ErrorInvalidStatusPlantao,
		plantao.ErrorInvalidTransitionStatus,
		plantao.ErrorValorTotalInvalido,
		plantao.ErrorOperacaoPagamentoObrigatoria,
		shared.ErrorEndBeforeStart,
		shared.ErrorPeriodoInvalido,
		usuario.ErrorInvalidEmail,
		usuario.ErrorPasswordShort,
		usuario.ErrorInvalidRole,
		setor.ErrorSetorNomeInvalido,
		smtp.ErrorSMTPHostInvalido,
		smtp.ErrorSMTPPortaInvalida,
		smtp.ErrorSMTPUsuarioInvalido,
		smtp.ErrorSMTPPasswordInvalido,
		smtp.ErrorSMTPFromInvalido,
		usuario.ErrorInvalidStatus,
	) {
		return http.StatusUnprocessableEntity, ErrorResponse{Code: "VALIDATION_ERROR", Message: err.Error()}
	}

	// 401
	if errors.Is(err, usuario.ErrInvalidCredentials) {
		return http.StatusUnauthorized, ErrorResponse{Code: "UNAUTHORIZED", Message: err.Error()}
	}

	if errors.Is(err, plantao.ErrorUsuarioSemPermissao) {
		return http.StatusForbidden, ErrorResponse{Code: "FORBIDDEN", Message: err.Error()}
	}

	// 500 — erro interno inesperado
	return http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_ERROR", Message: "erro interno do servidor"}
}

func isAny(target error, candidates ...error) bool {
	for _, c := range candidates {
		if errors.Is(target, c) {
			return true
		}
	}
	return false
}
