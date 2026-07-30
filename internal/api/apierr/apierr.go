package apierr

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"plantao/internal/domain/cargo"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/convite"
	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/setor"
	"plantao/internal/domain/shared"
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

func classify(err error) (int, ErrorResponse) {
	if errors.Is(err, financeiro.ErrorAtualizacaoVazia) {
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
		usuario.ErrorUserNotFound,
	) {
		return http.StatusNotFound, ErrorResponse{Code: "NOT_FOUND", Message: err.Error()}
	}

	// 409 — conflito
	if isAny(err,
		cargo.ErrorCargoAlreadyExists,
		colaborador.ErrorEmailAlreadyExists,
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
		financeiro.ErrorValorDiaNaoVigente,
		financeiro.ErrorConflitoValorDia,
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
		financeiro.ErrorValorDiaForaVigencia,
		financeiro.ErrorVigenciaInvalida,
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
