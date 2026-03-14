# Colaboradores

**Base URL:** `http://{host}/api/v1`

**Autenticação:** Bearer Token via header `Authorization: Bearer <token>`

> Todos os endpoints exigem autenticação (roles: `admin`, `gerente`, `colaborador`)

---

**Valores válidos para `cargo`:**

| Valor |
|-------|
| `Analista` |
| `Gerente` |
| `Consultor` |
| `Técnico` |
| `Outro` |
| `Desenvolvedor Frontend` |
| `Desenvolvedor Backend` |
| `Desenvolvedor Fullstack` |

**Valores válidos para `setor`:**

| Valor |
|-------|
| `TI` |
| `RH` |
| `Financeiro` |
| `Suporte` |
| `Desenvolvimento` |
| `Diretoria` |

---

## `POST /colaboradores`

> `Content-Type: application/json`

**Request:**
```json
{
  "nome": "João Silva",
  "email": "joao@email.com",
  "telefone": "(73) 99999-9999",
  "cargo": "Técnico",
  "setor": "TI",
  "status": "ativo",
  "ativo_plantao": "ativo",
  "data_admissao": "01/01/2024",
  "data_desligamento": null
}
```

> Datas no formato `DD/MM/AAAA` ou `AAAA-MM-DD`

**Response `201`:**
```json
{
  "id": "uuid",
  "nome": "João Silva",
  "email": "joao@email.com",
  "telefone": "(73) 99999-9999",
  "cargo": "Técnico",
  "setor": "TI",
  "foto_url": "",
  "status": "ativo",
  "ativo_plantao": "ativo",
  "data_admissao": "01/01/2024",
  "data_desligamento": ""
}
```

---

## `GET /colaboradores`

> Suporta filtros via query string

**Query params (todos opcionais):**
```
?nome=João&email=joao@&telefone=73&cargo=Técnico&setor=TI&data_admissao=01/01/2024
```

**Response `200`:** array de colaborador

---

## `GET /colaboradores/:id`

**Response `200`:** colaborador

---

## `PATCH /colaboradores/:id`

> `Content-Type: application/json` — todos os campos opcionais

**Request:**
```json
{
  "nome": "João Silva Atualizado",
  "cargo": "Analista",
  "setor": "Financeiro"
}
```

**Response `200`** (sem body)

---

## `PATCH /colaboradores/:id/foto`

> `Content-Type: multipart/form-data`

**Form field:**
| Campo | Tipo | Descrição |
|-------|------|-----------|
| `foto` | file | Imagem do colaborador |

**Response `200`** (sem body)

---

## `DELETE /colaboradores/:id`

> Desativa (soft delete) o colaborador

**Response `204`** (sem body)
