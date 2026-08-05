# Config Valores Dia

**Base URL:** `http://{host}/api/v1`

**Autenticação:** cookie `access_token`

> Todos os endpoints exigem autenticação e role `admin`.

**Tipos de dia:** `UTIL`, `SABADO`, `DOMINGO`, `FERIADO`

Cada tipo possui uma única configuração global. Não há vigência no contrato da
API: o mesmo valor é usado em qualquer data que ainda precise ser calculada.
Valores já registrados em `plantoes_detalhes` permanecem históricos e têm
prioridade nos relatórios.

As colunas de vigência continuam no banco somente por compatibilidade estrutural.
Novas linhas recebem internamente `vigencia_inicio = 1900-01-01` e
`vigencia_fim = NULL`; nenhuma leitura ou cálculo depende dessas datas.

---

## `GET /admin/config-valores`

Retorna todas as configurações globais, ordenadas por tipo.

**Response `200`:**

```json
[
  {
    "id": "uuid",
    "tipo_dia": "UTIL",
    "valor": 350.00
  }
]
```

---

## `POST /admin/config-valores`

Cria a configuração global de um tipo de dia.

**Request:**

```json
{
  "tipo_dia": "FERIADO",
  "valor": 500.00,
  "descricao": "Valor global para feriados"
}
```

**Response `201`:**

```json
{
  "id": "uuid",
  "tipo_dia": "FERIADO",
  "valor": 500.00,
  "descricao": "Valor global para feriados"
}
```

`descricao` é opcional. Se o tipo já estiver configurado, retorna `409 Conflict`; use o PATCH para
alterá-lo. `vigencia_inicio` e `vigencia_fim` não são aceitos e retornam `400`.

---

## `PATCH /admin/config-valores/:tipo_dia`

Atualiza parcialmente a configuração global. Campos omitidos são preservados e
`descricao: null` remove a descrição.

**Request:**

```json
{
  "valor": 175.50,
  "descricao": "Valor para dias úteis"
}
```

**Response `200`:**

```json
{
  "id": "uuid",
  "tipo_dia": "UTIL",
  "valor": 175.50,
  "descricao": "Valor para dias úteis",
  "updated_at": "2026-08-05T12:00:00Z"
}
```

Regras:

- somente `valor` e `descricao` são aceitos;
- `vigencia_inicio` e `vigencia_fim` retornam `400 Bad Request`;
- o valor deve ser positivo e possuir no máximo duas casas decimais;
- corpo vazio retorna `400 Bad Request`;
- tipo inexistente retorna `404 Not Found`;
- atualização concorrente retorna `409 Conflict`;
- gerente e colaborador recebem `403 Forbidden`.

Também existe `PATCH /admin/config-valores` exclusivamente para responder `400`
quando o tipo não for informado na URL.
