# Plantões

**Base URL:** `http://{host}/api/v1`

**Autenticação:** cookie `access_token`

> Todos os endpoints exigem autenticação (roles: `admin`, `gerente`, `colaborador`)

**Status do plantão:**
| Valor | Descrição |
|-------|-----------|
| `0`   | Agendado  |
| `1`   | Em andamento |
| `2`   | Concluído |
| `3`   | Cancelado |
| `4`   | Pago |

> Ao mover de `1` para `2`, o fechamento financeiro é executado em uma única
> transação. O status `4` não pode ser informado diretamente: ele é aplicado
> somente pelo endpoint de pagamento.

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

### Regras do fechamento

- `colaborador` pode fechar apenas o próprio plantão;
- `gerente` e `admin` podem fechar qualquer plantão;
- o período usa as datas civis de `America/Sao_Paulo`, incluindo início e fim;
- feriado prevalece sobre sábado ou domingo;
- cada data usa a configuração cuja vigência cobre aquela data;
- são criados detalhes, histórico e um pagamento `pendente`;
- repetição, dados preexistentes ou processamento concorrente retornam conflito,
  sem gravação parcial.

> A constraint atual `UNIQUE(tipo_dia)` permite somente uma linha de configuração
> por tipo de dia. Assim, sem alteração estrutural, o fechamento valida cada data
> contra essa única linha e falha caso ela esteja fora da vigência.

---

## `POST /plantoes/:id/pagamento`

Confirma o pagamento pendente criado no fechamento.

**Roles:** `admin`, `gerente`

**Request opcional:**

```json
{
  "observacoes": "Pagamento confirmado."
}
```

**Response `204`** (sem body)

O pagamento exige plantão concluído, detalhes consistentes e exatamente um
pagamento pendente com o mesmo colaborador e valor. A operação atualiza o
pagamento, muda o plantão para status `4` e grava o histórico atomicamente.
Repetições e chamadas concorrentes retornam `409 Conflict`.

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
