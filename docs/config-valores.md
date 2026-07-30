# Config Valores Dia

**Base URL:** `http://{host}/api/v1`

**Autenticação:** cookie `access_token`

> Todos os endpoints exigem autenticação (role: `admin`)

**Tipos de dia:** `UTIL`, `SABADO`, `DOMINGO`, `FERIADO`

**Sistema de vigência:** cada novo valor fecha automaticamente o anterior no dia anterior à nova `vigencia_inicio`.

> A estrutura atual possui `UNIQUE(tipo_dia)`, portanto aceita somente uma linha
> por tipo. O fluxo existente de criação de uma nova vigência tenta encerrar a
> linha anterior e inserir outra do mesmo tipo, o que não permite histórico real
> enquanto essa restrição existir. O endpoint de atualização abaixo modifica a
> linha vigente sem alterar o schema.

---

## `GET /admin/config-valores`
> Retorna os valores atualmente vigentes (um por tipo de dia)

**Response `200`:**
```json
[
  {
    "id": "uuid",
    "tipo_dia": "UTIL",
    "valor": 350.00,
    "vigencia_inicio": "2026-01-01T00:00:00Z",
    "vigencia_fim": null
  },
  {
    "id": "uuid",
    "tipo_dia": "SABADO",
    "valor": 450.00,
    "vigencia_inicio": "2026-01-01T00:00:00Z",
    "vigencia_fim": null
  },
  {
    "id": "uuid",
    "tipo_dia": "DOMINGO",
    "valor": 500.00,
    "vigencia_inicio": "2026-01-01T00:00:00Z",
    "vigencia_fim": null
  },
  {
    "id": "uuid",
    "tipo_dia": "FERIADO",
    "valor": 600.00,
    "vigencia_inicio": "2026-01-01T00:00:00Z",
    "vigencia_fim": null
  }
]
```

---

## `PATCH /admin/config-valores/:tipo_dia`

Atualiza parcialmente a configuração que está vigente na data civil atual de
`America/Sao_Paulo`. A rota aceita exclusivamente a role `admin`.

O `tipo_dia` deve ser informado na URL usando um dos valores `UTIL`, `SABADO`,
`DOMINGO` ou `FERIADO`. Campos omitidos são preservados. Em campos anuláveis,
`null` remove o valor.

**Request:**

```json
{
  "valor": 175.50,
  "descricao": "Valor para dias úteis",
  "vigencia_inicio": "2026-05-01",
  "vigencia_fim": null
}
```

**Response `200`:**

```json
{
  "id": "uuid",
  "tipo_dia": "UTIL",
  "valor": 175.50,
  "descricao": "Valor para dias úteis",
  "vigencia_inicio": "2026-05-01T00:00:00Z",
  "vigencia_fim": null,
  "updated_at": "2026-07-30T12:00:00Z"
}
```

Regras:

- a configuração precisa estar vigente antes e depois da alteração;
- `vigencia_fim` não pode ser anterior a `vigencia_inicio`;
- o valor deve ser positivo e possuir no máximo duas casas decimais;
- requisições concorrentes para o mesmo tipo retornam `409 Conflict`;
- um corpo vazio ou a ausência de `tipo_dia` retorna `400 Bad Request`;
- gerente e colaborador recebem `403 Forbidden`.

---

## `POST /admin/config-valores`
> Define um novo valor para um tipo de dia. Se já existir um valor vigente para o mesmo tipo, ele é encerrado automaticamente no dia anterior à nova vigência.

**Request:**
```json
{
  "tipo_dia": "UTIL",
  "valor": 400.00,
  "vigencia_inicio": "2026-06-01"
}
```

**Response `201`:**
```json
{
  "id": "uuid",
  "tipo_dia": "UTIL",
  "valor": 400.00,
  "vigencia_inicio": "2026-06-01T00:00:00Z",
  "vigencia_fim": null
}
```
