# Configuração SMTP

**Base URL:** `http://{host}/api/v1`

Todos os endpoints exigem autenticação e role `admin`.

## `GET /admin/config-smtp`

Retorna a configuração SMTP global. A senha nunca é retornada. Quando não há
configuração cadastrada, responde `404`.

## `POST /admin/config-smtp`

Cria a única configuração SMTP do sistema.

```json
{
  "smtp_host": "smtp.example.com",
  "smtp_port": 587,
  "smtp_username": "usuario",
  "smtp_password": "senha",
  "smtp_from": "no-reply@example.com"
}
```

Todos os campos são obrigatórios. Uma segunda criação responde `409 Conflict`.

## `PATCH /admin/config-smtp`

Atualiza parcialmente a configuração. Campos omitidos são preservados;
`smtp_password` omitida mantém a senha atual. Valores `null` não são aceitos.
As respostas de `POST`, `GET` e `PATCH` omitem `smtp_password`.
