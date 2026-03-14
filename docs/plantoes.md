# Plantões

**Base URL:** `http://{host}/api/v1`

**Autenticação:** Bearer Token via header `Authorization: Bearer <token>`

> Todos os endpoints exigem autenticação (roles: `admin`, `gerente`, `colaborador`)

**Status do plantão:**
| Valor | Descrição |
|-------|-----------|
| `0`   | Pendente  |
| `1`   | Em andamento |
| `2`   | Concluído |

> Ao mover para status `2` (Concluído), o cálculo financeiro é disparado automaticamente.

---

## `POST /plantoes`

**Request:**
```json
{
  "colaborador_id": "uuid",
  "periodo": {
    "inicio": "2026-01-01",
    "fim": "2026-01-07"
  }
}
```

**Response `201`:**
```json
{
  "id": "uuid",
  "colaborador_id": "uuid",
  "periodo": {
    "inicio": "2026-01-01T00:00:00Z",
    "fim": "2026-01-07T00:00:00Z"
  },
  "status": 0,
  "valor_total": 0
}
```

---

## `GET /plantoes`
> Retorna todos os plantões

**Response `200`:** array de plantão

---

## `GET /plantoes/:id`

**Response `200`:** plantão

---

## `DELETE /plantoes/:id`

**Response `204`** (sem body)

---

## `PATCH /plantoes/:id/status`

**Request:**
```json
{
  "new_status": "2",
  "observacoes": "Plantão encerrado sem intercorrências."
}
```

> `observacoes` é opcional e só é persistido quando `new_status` for `2` (Concluído).

**Response `204`** (sem body)

---

## `GET /plantoes/colaborador/:colaborador_id`

**Response `200`:** array de plantão do colaborador

---

## `GET /plantoes/status/:status`

**Params:** `:status` = `0`, `1` ou `2`

**Response `200`:** array de plantão

---

## `GET /plantoes/periodo/:start_date/:end_date`

**Params:** datas no formato `YYYY-MM-DD`

**Exemplo:** `/plantoes/periodo/2026-01-01/2026-01-31`

**Response `200`:** array de plantão