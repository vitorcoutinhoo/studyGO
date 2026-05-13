# API Endpoints — studyGO

**Base URL:** `http://localhost:{PORT}/api/v1`

**Autenticação:** Bearer Token (JWT) via header `Authorization: Bearer <token>`

**Roles disponíveis:** `ADMIN`, `GERENTE`, `COLABORADOR`

---

## Sumário

- [Auth](#auth)
- [Usuários](#usuários)
- [Colaboradores](#colaboradores)
- [Plantões](#plantões)
- [Convites](#convites)
- [Feriados](#feriados)
- [Config Valores](#config-valores)
- [Modelos de Comunicação](#modelos-de-comunicação)
- [Cargos](#cargos)
- [Setores](#setores)
- [Roles](#roles)

---

## Auth

### POST `/api/v1/auth/login`
Login de usuário. Sujeito a rate limiting.

**Request:**
```json
{
  "email": "usuario@email.com",
  "senha": "minimo6chars"
}
```

**Response `200`:**
```json
{
  "id": "uuid",
  "id_colaborador": "uuid",
  "email": "usuario@email.com",
  "role": "COLABORADOR",
  "ativo": "S"
}
```

---

## Usuários

### POST `/api/v1/usuarios/cadastro`
Cria um usuário a partir de um token de convite. Público.

**Request:**
```json
{
  "email": "usuario@email.com",
  "senha": "minimo6chars"
}
```
> O token de convite deve ser enviado como query param `?token=<token>` ou conforme implementado no frontend.

**Response `201`:**
```json
{
  "id": "uuid",
  "id_colaborador": "uuid",
  "email": "usuario@email.com",
  "role": "COLABORADOR",
  "ativo": "S"
}
```

---

### GET `/api/v1/authenticated/usuarios`
Retorna os dados do usuário autenticado.

**Auth:** `COLABORADOR`, `GERENTE`, `ADMIN`

**Response `200`:**
```json
{
  "id": "uuid",
  "id_colaborador": "uuid",
  "email": "usuario@email.com",
  "role": "COLABORADOR",
  "ativo": "S"
}
```

---

### PUT `/api/v1/authenticated/usuarios`
Atualiza email e/ou senha do usuário autenticado.

**Auth:** `COLABORADOR`, `GERENTE`, `ADMIN`

**Request:**
```json
{
  "email": "novo@email.com",
  "senha": "novasenha123"
}
```

**Response `200`:** sem body

---

### DELETE `/api/v1/authenticated/usuarios`
Desativa a conta do usuário autenticado.

**Auth:** `COLABORADOR`, `GERENTE`, `ADMIN`

**Response `204`:** sem body

---

### GET `/api/v1/admin/all`
Lista todos os usuários do sistema.

**Auth:** `ADMIN`

**Response `200`:**
```json
[
  {
    "id": "uuid",
    "id_colaborador": "uuid",
    "email": "usuario@email.com",
    "role": "COLABORADOR",
    "ativo": "S"
  }
]
```

---

## Colaboradores

### POST `/api/v1/colaboradores`
Cria um novo colaborador.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Request:**
```json
{
  "nome": "João Silva",
  "email": "joao@email.com",
  "telefone": "(11) 99999-0000",
  "cargo": "Analista",
  "setor": "TI",
  "status": "ATIVO",
  "ativo_plantao": "S",
  "data_admissao": "2024-01-15",
  "data_desligamento": null
}
```

**Response `201`:**
```json
{
  "id": "uuid",
  "nome": "João Silva",
  "email": "joao@email.com",
  "telefone": "(11) 99999-0000",
  "cargo": "Analista",
  "setor": "TI",
  "foto_url": "https://...",
  "status": "ATIVO",
  "ativo_plantao": "S",
  "data_admissao": "2024-01-15",
  "data_desligamento": ""
}
```

---

### GET `/api/v1/colaboradores`
Lista colaboradores com filtros opcionais via query params.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Query params (todos opcionais):**
| Param | Tipo | Descrição |
|---|---|---|
| `nome` | string | Filtrar por nome |
| `email` | string | Filtrar por email |
| `telefone` | string | Filtrar por telefone |
| `cargo` | string | Filtrar por cargo |
| `setor` | string | Filtrar por setor |
| `data_admissao` | string | Filtrar por data de admissão |

**Response `200`:**
```json
[
  {
    "id": "uuid",
    "nome": "João Silva",
    "email": "joao@email.com",
    "telefone": "(11) 99999-0000",
    "cargo": "Analista",
    "setor": "TI",
    "foto_url": "https://...",
    "status": "ATIVO",
    "ativo_plantao": "S",
    "data_admissao": "2024-01-15",
    "data_desligamento": ""
  }
]
```

---

### GET `/api/v1/colaboradores/:id`
Busca um colaborador pelo ID.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `200`:** mesmo objeto acima

---

### PATCH `/api/v1/colaboradores/:id`
Atualiza dados de um colaborador (somente campos enviados são alterados).

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Request (todos os campos são opcionais):**
```json
{
  "nome": "João Atualizado",
  "email": "novo@email.com",
  "telefone": "(11) 88888-0000",
  "cargo": "Sênior",
  "setor": "TI",
  "status": "ATIVO",
  "ativo_plantao": "N",
  "data_admissao": "2024-01-15",
  "data_desligamento": "2025-12-31"
}
```

**Response `200`:** sem body

---

### PATCH `/api/v1/colaboradores/:id/foto`
Faz upload da foto do colaborador. Limite: 5MB.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Request:** `multipart/form-data`
| Campo | Tipo | Descrição |
|---|---|---|
| `foto` | file | Imagem (JPEG, PNG, etc.) |

**Response `200`:** sem body

---

### DELETE `/api/v1/colaboradores/:id`
Desativa um colaborador (soft delete).

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `204`:** sem body

---

## Plantões

### POST `/api/v1/plantoes`
Cria um novo plantão.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Request:**
```json
{
  "colaborador_id": "uuid",
  "periodo": {
    "inicio": "2025-05-01T08:00:00Z",
    "fim": "2025-05-01T20:00:00Z"
  }
}
```

**Response `201`:**
```json
{
  "id": "uuid",
  "colaborador_id": "uuid",
  "periodo": {
    "inicio": "2025-05-01T08:00:00Z",
    "fim": "2025-05-01T20:00:00Z"
  },
  "status": "PENDENTE",
  "valor_total": 250.00,
  "observacoes": null
}
```

---

### GET `/api/v1/plantoes`
Lista todos os plantões.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `200`:** array do objeto acima

---

### GET `/api/v1/plantoes/:id`
Busca um plantão pelo ID.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `200`:** objeto acima

---

### GET `/api/v1/plantoes/colaborador/:colaborador_id`
Lista plantões de um colaborador específico.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `200`:** array do objeto acima

---

### GET `/api/v1/plantoes/status/:status`
Lista plantões por status.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Valores de `:status`:** `PENDENTE`, `APROVADO`, `RECUSADO` *(confirme com o backend)*

**Response `200`:** array do objeto acima

---

### GET `/api/v1/plantoes/periodo/:start_date/:end_date`
Lista plantões em um intervalo de datas.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Formato das datas:** `YYYY-MM-DD`

**Exemplo:** `/api/v1/plantoes/periodo/2025-05-01/2025-05-31`

**Response `200`:** array do objeto acima

---

### PATCH `/api/v1/plantoes/:id/status`
Atualiza o status de um plantão.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Request:**
```json
{
  "new_status": "APROVADO",
  "observacoes": "Aprovado pelo gerente"
}
```

**Response `204`:** sem body

---

### DELETE `/api/v1/plantoes/:id`
Remove um plantão.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `204`:** sem body

---

## Convites

### POST `/api/v1/convites`
Cria um convite de cadastro para um colaborador.

**Auth:** `ADMIN`

**Request:**
```json
{
  "id_colaborador": "uuid",
  "colaborador_email": "joao@email.com",
  "colaborador_nome": "João Silva"
}
```

**Response `201`:** sem body (envio de email disparado automaticamente)

---

### GET `/api/v1/convites`
Lista todos os convites.

**Auth:** `ADMIN`

**Response `200`:**
```json
[
  {
    "token": "string",
    "id_colaborador": "uuid",
    "expira_em": "2025-05-10T00:00:00Z",
    "usado": false,
    "created_at": "2025-05-01T10:00:00Z"
  }
]
```

---

### DELETE `/api/v1/convites/:token`
Desativa um convite pelo token.

**Auth:** `ADMIN`

**Response `204`:** sem body

---

## Feriados

### GET `/api/v1/admin/feriados`
Lista os feriados de um ano.

**Auth:** `ADMIN`

**Query params:**
| Param | Tipo | Descrição |
|---|---|---|
| `ano` | int | Ano desejado. Padrão: ano atual |

**Response `200`:**
```json
[
  {
    "id": "uuid",
    "data": "2025-01-01T00:00:00Z",
    "nome": "Confraternização Universal",
    "descricao": "Ano Novo"
  }
]
```

---

### PATCH `/api/v1/admin/feriados/:id/data`
Atualiza a data de um feriado.

**Auth:** `ADMIN`

**Request:**
```json
{
  "nova_data": "2025-11-15"
}
```

**Response `200`:**
```json
{
  "id": "uuid",
  "data": "2025-11-15T00:00:00Z",
  "nome": "Proclamação da República",
  "descricao": ""
}
```

---

## Config Valores

Configuração dos valores pagos por tipo de dia de plantão.

### GET `/api/v1/admin/config-valores`
Lista as configurações de valor vigentes.

**Auth:** `ADMIN`

**Response `200`:**
```json
[
  {
    "id": "uuid",
    "tipo_dia": "UTIL",
    "valor": 250.00,
    "vigencia_inicio": "2025-01-01T00:00:00Z",
    "vigencia_fim": null
  }
]
```

---

### POST `/api/v1/admin/config-valores`
Define um novo valor para um tipo de dia.

**Auth:** `ADMIN`

**Request:**
```json
{
  "tipo_dia": "FERIADO",
  "valor": 500.00,
  "vigencia_inicio": "2025-06-01"
}
```

**Valores de `tipo_dia`:** `UTIL`, `SABADO`, `DOMINGO`, `FERIADO` *(confirme com o backend)*

**Response `201`:**
```json
{
  "id": "uuid",
  "tipo_dia": "FERIADO",
  "valor": 500.00,
  "vigencia_inicio": "2025-06-01T00:00:00Z",
  "vigencia_fim": null
}
```

---

## Modelos de Comunicação

Templates de email/comunicação gerenciados pelo admin.

### POST `/api/v1/auth/admin/modelo-comunicacao/`
Cria um novo modelo.

**Auth:** `ADMIN`

**Request:**
```json
{
  "nome": "Boas-vindas",
  "tipo_comunicacao": "EMAIL",
  "assunto": "Bem-vindo ao sistema",
  "corpo": "Olá {{nome}}, seu acesso foi criado."
}
```

**Response `201`:**
```json
{
  "id": "uuid",
  "nome": "Boas-vindas",
  "tipo_comunicacao": "EMAIL",
  "assunto": "Bem-vindo ao sistema",
  "corpo": "Olá {{nome}}, seu acesso foi criado.",
  "ativo": "S"
}
```

---

### GET `/api/v1/auth/admin/modelo-comunicacao/`
Lista todos os modelos.

**Auth:** `ADMIN`

**Response `200`:** array do objeto acima

---

### GET `/api/v1/auth/admin/modelo-comunicacao/:id_modelo`
Busca um modelo pelo ID.

**Auth:** `ADMIN`

**Response `200`:** objeto acima

---

### PUT `/api/v1/auth/admin/modelo-comunicacao/:id_modelo`
Atualiza um modelo.

**Auth:** `ADMIN`

**Request (todos opcionais):**
```json
{
  "nome": "Boas-vindas v2",
  "tipo_comunicacao": "EMAIL",
  "assunto": "Bem-vindo!",
  "corpo": "Olá {{nome}}, seja bem-vindo.",
  "ativo": "S"
}
```

**Response `200`:** sem body

---

### DELETE `/api/v1/auth/admin/modelo-comunicacao/:id_modelo`
Desativa um modelo.

**Auth:** `ADMIN`

**Response `204`:** sem body

---

## Cargos

### GET `/api/v1/cargos`
Lista todos os cargos ativos.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `200`:**
```json
[
  { "id": "uuid", "nome": "Desenvolvedor Backend" }
]
```

---

### POST `/api/v1/admin/cargos`
Cria um novo cargo.

**Auth:** `ADMIN`

**Request:**
```json
{ "nome": "Arquiteto de Software" }
```

**Response `201`:**
```json
{ "id": "uuid", "nome": "Arquiteto de Software" }
```

---

### DELETE `/api/v1/admin/cargos/:id`
Desativa um cargo.

**Auth:** `ADMIN`

**Response `204`:** sem body

---

## Setores

### GET `/api/v1/setores`
Lista todos os setores ativos.

**Auth:** `ADMIN`, `GERENTE`, `COLABORADOR`

**Response `200`:**
```json
[
  { "id": "uuid", "nome": "TI" }
]
```

---

### POST `/api/v1/admin/setores`
Cria um novo setor.

**Auth:** `ADMIN`

**Request:**
```json
{ "nome": "Jurídico" }
```

**Response `201`:**
```json
{ "id": "uuid", "nome": "Jurídico" }
```

---

### DELETE `/api/v1/admin/setores/:id`
Desativa um setor.

**Auth:** `ADMIN`

**Response `204`:** sem body

---

## Roles

### GET `/api/v1/roles`
Lista todas as roles disponíveis. Público.

**Response `200`:**
```json
[
  { "id": "uuid", "nome": "admin" },
  { "id": "uuid", "nome": "gerente" },
  { "id": "uuid", "nome": "colaborador" }
]
```

---

## Códigos de Erro

Todos os erros seguem o formato:

```json
{
  "code": "NOT_FOUND",
  "message": "mensagem descritiva do erro"
}
```

| Status | Significado |
|---|---|
| `400` | Bad Request — dados inválidos ou campos faltando |
| `401` | Unauthorized — token ausente ou inválido |
| `403` | Forbidden — sem permissão para este recurso |
| `404` | Not Found — recurso não encontrado |
| `409` | Conflict — conflito (ex: email já cadastrado) |
| `422` | Unprocessable Entity — regra de negócio violada |
| `500` | Internal Server Error |
