# Feriados

**Base URL:** `http://{host}/api/v1`

**Autenticação:** Bearer Token via header `Authorization: Bearer <token>`

> Todos os endpoints exigem autenticação (role: `admin`)

**Tipos de feriado (campo `descricao`):** `NACIONAL`, `ESTADUAL`, `MUNICIPAL`, `COMEMORATIVO`, `FACULTATIVO`

---

## `GET /admin/feriados?ano=2026`

**Query param obrigatório:** `ano` (ex: `2026`)

**Response `200`:**
```json
[
  {
    "id": "uuid",
    "data": "2026-01-01T00:00:00Z",
    "nome": "Confraternização Universal",
    "descricao": "NACIONAL"
  }
]
```

---

## `PATCH /admin/feriados/:id/data`
> Altera a data de um feriado **municipal**. Retorna erro se o feriado não for do tipo `MUNICIPAL`.

**Request:**
```json
{
  "nova_data": "2026-11-20"
}
```

**Response `200`:** feriado atualizado
