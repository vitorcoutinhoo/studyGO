# Colaboradores

**Base URL:** `http://{host}/api/v1`

**Autenticação:** Bearer Token via header `Authorization: Bearer <token>`

> Todos os endpoints exigem autenticação (roles: `admin`, `gerente`, `colaborador`)

---

## `POST /colaboradores`

**Request:**
```json
{
  "nome": "João Silva",
  "email": "joao@email.com",
  "telefone": "73999999999",
  "cargo": "MEDICO",
  "setor": "URGENCIA",
  "foto_url": "https://...",
  "status": "ativo",
  "ativo_plantao": "ativo",
  "data_admissao": "01/01/2024",
  "data_desligamento": null
}
```
> Datas no formato `DD/MM/YYYY`

**Response `201`:**
```json
{
  "id": "uuid",
  "nome": "João Silva",
  "email": "joao@email.com",
  "telefone": "73999999999",
  "cargo": "MEDICO",
  "setor": "URGENCIA",
  "foto_url": "https://...",
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
?nome=João&email=joao@&telefone=73&cargo=MEDICO&setor=URGENCIA&data_admissao=01/01/2024
```

**Response `200`:** array de colaborador

---

## `GET /colaboradores/:id`

**Response `200`:** colaborador

---

## `PATCH /colaboradores/:id`

**Request:** mesmo schema do `POST`, todos os campos opcionais

**Response `200`** (sem body)

---

## `DELETE /colaboradores/:id`
> Desativa (soft delete) o colaborador

**Response `204`** (sem body)