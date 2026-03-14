# Config Valores Dia

**Base URL:** `http://{host}/api/v1`

**Autenticação:** Bearer Token via header `Authorization: Bearer <token>`

> Todos os endpoints exigem autenticação (role: `admin`)

**Tipos de dia:** `UTIL`, `SABADO`, `DOMINGO`, `FERIADO`

**Sistema de vigência:** cada novo valor fecha automaticamente o anterior no dia anterior à nova `vigencia_inicio`.

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