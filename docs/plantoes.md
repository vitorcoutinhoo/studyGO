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

## Início automático

Ao iniciar a API, um worker interno procura imediatamente plantões `agendados`
cujo instante de início já foi alcançado e os move para `em andamento`. Depois
disso, a verificação é repetida a cada 5 minutos; portanto, em operação normal,
a transição pode ocorrer com atraso aproximado de até 5 minutos.

Status e histórico são gravados na mesma transação. O histórico identifica a
alteração como automática e não possui usuário (`id_usuario = NULL`). Os lotes
são pequenos e usam bloqueio com `SKIP LOCKED`, permitindo múltiplas instâncias
da API sem processar o mesmo plantão simultaneamente. Plantões ignorados por
bloqueios concorrentes ou acumulados durante uma indisponibilidade são
reconsiderados nos ciclos seguintes.

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
- cada data usa a configuração global do seu tipo de dia;
- são criados detalhes, histórico e um pagamento `pendente`;
- repetição, dados preexistentes ou processamento concorrente retornam conflito,
  sem gravação parcial.

> A constraint `UNIQUE(tipo_dia)` garante uma única configuração global para cada
> tipo. As datas internas de vigência não participam do fechamento.

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

## `GET /plantoes/relatorio`

Retorna uma linha para cada dia civil de plantão. A consulta exige autenticação
e aceita somente as roles `admin` e `gerente`.

**Query parameters:**

- `data_inicio`: obrigatório, no formato `YYYY-MM-DD`;
- `data_fim`: obrigatório, no formato `YYYY-MM-DD`;
- `colaborador_id`: UUID opcional; quando omitido, inclui todos os colaboradores;
- `status`: opcional, com valores de `0` a `4`; quando omitido, inclui qualquer status.

**Exemplo:**

```http
GET /api/v1/plantoes/relatorio?data_inicio=2026-08-01&data_fim=2026-08-31&status=2
```

**Response `200`:**

```json
[
  {
    "colaborador_id": "uuid",
    "plantao_id": "uuid",
    "status": 2,
    "data_inicio": "2026-08-01T03:00:00Z",
    "data_fim": "2026-08-03T03:00:00Z",
    "data": "2026-08-01",
    "nome_colaborador": "Maria Silva",
    "valor_total": 450.00,
    "valor": 150.00,
    "observacoes": null
  }
]
```

O intervalo é inclusivo e aplicado a cada data civil entre o início e o fim do
plantão em `America/Sao_Paulo`. Quando existir detalhe persistido para a data,
seu valor histórico será utilizado. Caso contrário, o relatório determina se o
dia é útil, sábado, domingo ou feriado e usa a configuração global do tipo.
Feriados têm prioridade sobre o dia da semana.

O campo `valor_total` usa o valor persistido no plantão quando ele for diferente
de zero. Quando o valor persistido for `0` ou `NULL`, o relatório soma todos os
dias do plantão completo, ainda que o período consultado mostre somente parte
deles. Esse total é calculado apenas para a resposta e não é gravado no banco.

Os dias calculados existem somente na resposta: consultar o relatório não grava
linhas em `plantoes_detalhes`. Assim, plantões ainda sem detalhes também aparecem.
Uma consulta sem resultados responde com `[]`.

Parâmetros ausentes ou malformados retornam `400 Bad Request`. Período invertido
ou status numérico fora de `0` a `4` retorna `422 Unprocessable Entity`. Um dia
sem detalhe exige uma configuração global do tipo correspondente. Quando o total
precisar ser calculado, essa validação cobre todos os dias do plantão completo;
configuração ausente retorna `404`.

---

## `GET /plantoes/periodo/:start_date/:end_date`

**Params:** datas no formato `YYYY-MM-DD`

**Exemplo:** `/plantoes/periodo/2026-01-01/2026-01-31`

**Response `200`:** array de plantão
