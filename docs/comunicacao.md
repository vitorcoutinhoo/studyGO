# Comunicação — Modelos

**Base URL:** `http://{host}/api/v1`

**Autenticação:** Bearer Token via header `Authorization: Bearer <token>`

> Todos os endpoints exigem autenticação (role: `admin`)

---

**Valores válidos para `tipo_comunicacao`:**

| Valor | Tags obrigatórias no corpo |
|-------|---------------------------|
| `Plantão Agendado` | `{{.nome}}`, `{{.dataInicio}}`, `{{.dataFim}}` |
| `Plantão Concluido` | `{{.nome}}`, `{{.dataInicio}}`, `{{.dataFim}}` |
| `Plantão Ainda Está Aberto` | `{{.nome}}`, `{{.dataInicio}}`, `{{.dataFim}}` |
| `Plantão Pago` | `{{.nome}}`, `{{.dataInicio}}`, `{{.dataFim}}`, `{{.valorPago}}` |
| `Colaborador Cadastrado` | `{{.nome}}`, `{{.email}}` |
| `Colaborador Atualizado` | `{{.nome}}`, `{{.email}}` |
| `Colaborador Deletado` | `{{.nome}}`, `{{.dataAtual}}` |
| `Usuário Cadastrado` | `{{.nome}}`, `{{.email}}` |
| `Email do Usuário Atualizado` | `{{.nome}}`, `{{.email}}` |
| `Senha do Usuário Atualizada` | `{{.nome}}`, `{{.dataAtual}}` |
| `Usuário Deletado` | `{{.nome}}`, `{{.dataAtual}}` |

> O `corpo` deve ser HTML válido e conter todas as tags obrigatórias do tipo escolhido. Scripts não são permitidos.

---

## `POST /auth/admin/modelo-comunicacao/`

**Request:**
```json
{
  "nome": "Aviso de Plantão Agendado",
  "tipo_comunicacao": "Plantão Agendado",
  "assunto": "Seu plantão foi agendado",
  "corpo": "<p>Olá {{.nome}}, seu plantão foi agendado de {{.dataInicio}} até {{.dataFim}}.</p>"
}
```

**Response `201`:**
```json
{
  "id": "uuid",
  "nome": "Aviso de Plantão Agendado",
  "tipo_comunicacao": "Plantão Agendado",
  "assunto": "Seu plantão foi agendado",
  "corpo": "<p>...</p>",
  "ativo": "ativo"
}
```

---

## `GET /auth/admin/modelo-comunicacao/`

**Response `200`:** array de modelos

---

## `GET /auth/admin/modelo-comunicacao/:id_modelo`

**Response `200`:** modelo

---

## `PUT /auth/admin/modelo-comunicacao/:id_modelo`

**Request:**
```json
{
  "nome": "Novo nome",
  "tipo_comunicacao": "Plantão Pago",
  "assunto": "Novo assunto",
  "corpo": "<p>Olá {{.nome}}, seu plantão de {{.dataInicio}} a {{.dataFim}} foi pago. Valor: {{.valorPago}}</p>",
  "ativo": "ativo"
}
```

**Response `200`** (sem body)

---

## `DELETE /auth/admin/modelo-comunicacao/:id_modelo`
> Desativa (soft delete) o modelo

**Response `204`** (sem body)
