# Autenticação

**Base URL:** `http://{host}/api/v1`

---

## `POST /auth/login`
> Público — sem autenticação

**Request:**
```json
{
  "email": "usuario@email.com",
  "senha": "minhasenha"
}
```

**Response `200`:**
```json
{
  "token": "eyJhbGci..."
}
```