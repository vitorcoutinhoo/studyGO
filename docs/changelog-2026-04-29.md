# Changelog — 29/04/2026

## Resumo

Quatro commits realizados hoje com foco em qualidade de código, consistência de datas, padronização de erros HTTP e segurança de entrada de dados.

---

## [refactor] Struct de auditoria compartilhada e suporte a timezone nas datas

**Commit:** `04ec192`

### Motivação
Os campos `CreatedAt` e `UpdatedAt` estavam duplicados em cada domínio como `*time.Time`, sem convenção de timezone. As funções de conversão de data não suportavam fuso horário.

### Mudanças

**`internal/domain/shared/auditoria.go`** *(novo)*
- Criada struct `Auditoria` com campos `CreatedAt time.Time` e `UpdatedAt time.Time`

**Domínios atualizados para embedar `shared.Auditoria`**
- `internal/domain/usuario/usuario.go`
- `internal/domain/colaborador/colaborador.go`
- `internal/domain/plantao/plantao.go`
- `internal/domain/comunicacao/modelo.go`

**`internal/utils/date_converter.go`**
- `ParseBrToUsDate` e `ParseUsToBrDate` agora aceitam `*time.Location` como parâmetro
- Usa `time.ParseInLocation` para respeitar o fuso informado
- Adicionado `BrasilLocation()` — retorna `America/Sao_Paulo`
- Adicionado `ToLocation(t, loc)` — converte um `time.Time` para outro fuso

**`seed/01_schema.sql`**
- Todas as colunas `TIMESTAMP` convertidas para `TIMESTAMPTZ` (colunas de auditoria e timestamps de eventos)
- Colunas `DATE` mantidas onde apenas a data importa (admissão, desligamento, feriados, vigência)

**Arquivos de infraestrutura ajustados**
- `internal/infra/persistence/postgres/modelo_repository.go` — atribuições de `CreatedAt`/`UpdatedAt` corrigidas para `time.Time` (sem ponteiro)
- `internal/domain/usuario/auth.go` — inicialização via `Auditoria: user.Auditoria`

---

## [fix] Validação de email por regex RFC-compatível

**Commit:** `036d5aa`

### Motivação
A validação anterior (`len(email) <= 30 && strings.Contains(email, "@")`) aceitava strings inválidas como `a@`, `@@` ou `x@x`.

### Mudanças

**`internal/domain/colaborador/colaborador.go`**
**`internal/domain/usuario/usuario.go`**

- Substituído por `regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)`
- Limite de tamanho ampliado de 30 para 100 caracteres
- Import de `strings` removido; adicionado `regexp`

---

## [feat] Binding tags de validação em todos os DTOs de request

**Commit:** `901ac39`

### Motivação
A maioria dos DTOs não tinha tags `binding:`, permitindo que campos obrigatórios chegassem vazios ao domínio sem rejeição na camada HTTP.

### Mudanças por arquivo

**`internal/api/dto/colaboradorDto.go`**
| Campo | Tag adicionada |
|---|---|
| `Nome` | `binding:"required,min=2,max=100"` |
| `Email` | `binding:"required,email,max=100"` |
| `Telefone` | `binding:"required"` |
| `Cargo` | `binding:"required"` |
| `Setor` | `binding:"required"` |
| `DataAdmissao` | `binding:"required"` |
| `Nome` (Update) | `binding:"omitempty,min=2,max=100"` |
| `Email` (Update) | `binding:"omitempty,email,max=100"` |

**`internal/api/dto/usuarioDto.go`**
| Campo | Tag adicionada |
|---|---|
| `Email` | `binding:"required,email,max=100"` |
| `Senha` | `binding:"required,min=6,max=72"` |
| `Role` (Admin) | `binding:"required,oneof=admin gerente colaborador"` |

**`internal/api/dto/conviteDto.go`**
| Campo | Tag adicionada |
|---|---|
| `IdColaborador` | `binding:"required,uuid"` |
| `ColaboradorEmail` | `binding:"required,email,max=100"` |
| `ColaboradorNome` | `binding:"required,min=2,max=100"` |

**`internal/api/dto/modeloComunicacaoDTO.go`**
| Campo | Tag adicionada |
|---|---|
| `Nome` | `binding:"required,min=1,max=255"` |
| `TipoComunicacao` | `binding:"required"` |
| `Assunto` | `binding:"required,min=1,max=255"` |
| `Corpo` | `binding:"required,min=1"` |

---

## [feat] Erros HTTP estruturados, status codes corretos e limite de 5 MB no upload de foto

**Commit:** `aafb771`

### Motivação
Todos os erros de domínio retornavam `500 Internal Server Error` com `gin.H{"error": "..."}`. Isso dificultava o tratamento no frontend e escondia a natureza real do erro. O upload de foto não tinha limite de tamanho.

### Mudanças

**`internal/api/apierr/apierr.go`** *(novo)*

Pacote centralizado de erros HTTP com:

```go
type ErrorResponse struct {
    Code    string `json:"code"`
    Message string `json:"message"`
}

func Respond(ctx *gin.Context, err error)
```

Mapeamento de erros de domínio para status HTTP:

| Status | Code | Erros de domínio |
|---|---|---|
| `404` | `NOT_FOUND` | `ErrorColaboradorNotFound`, `ErrorPlantaoNotFinded`, `ErrorUserNotFound`, etc. |
| `409` | `CONFLICT` | `ErrorEmailAlreadyExists`, `ErrorExistingPlantao`, `ErrorModeloComunicacaoAlreadyExists` |
| `410` | `GONE` | `ErrorConviteExpired`, `ErrorConviteUsed` |
| `422` | `VALIDATION_ERROR` | Erros de regra de negócio (status inválido, cargo inválido, período inválido, etc.) |
| `401` | `UNAUTHORIZED` | `ErrInvalidCredentials` |
| `500` | `INTERNAL_ERROR` | Erros inesperados (mensagem genérica, sem vazar detalhes internos) |

**Controllers atualizados** (todos os 8):
- `auth_controller.go`
- `colaborador_controller.go`
- `convite_controller.go`
- `feriado_controller.go`
- `modelo_comunicacaoController.go`
- `plantao_controller.go`
- `usuario_controller.go`
- `valor_dia_controller.go`

**Upload de foto** (`colaborador_controller.go`):
- Constante `maxFotoSize = 5 MB`
- Verificação em 3 pontos: `CreateColaborador`, `UpdateColaborador`, `UploadFotoColaborador`
- Retorna `413 Request Entity Too Large` com `code: "FILE_TOO_LARGE"` se exceder o limite
- O storage já validava o tipo MIME (jpeg, png, webp) — mantido sem alteração

---

## Arquivos criados hoje

| Arquivo | Tipo |
|---|---|
| `internal/domain/shared/auditoria.go` | Novo |
| `internal/api/apierr/apierr.go` | Novo |
| `docs/changelog-2026-04-29.md` | Novo |
