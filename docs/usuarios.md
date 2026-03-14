# Usuários

**Base URL:** `http://{host}/api/v1`

**Autenticação:** Bearer Token via header `Authorization: Bearer <token>`

---

## `POST /usuarios/colaboradores/:id_colaborador`

> Público — cria o acesso de login para um colaborador já cadastrado

**Request:**
```json
{
  "email": "colaborador@email.com",
  "senha": "minhasenha"
}
```

**Response `201`:**
```json
{
  "id": "uuid",
  "id_colaborador": "uuid",
  "email": "colaborador@email.com",
  "role": "colaborador",
  "ativo": "ativo"
}
```

---

## `GET /authenticated/usuarios`

> Retorna os dados do usuário autenticado (roles: `colaborador`, `gerente`, `admin`)

**Response `200`:**
```json
{
  "id": "uuid",
  "id_colaborador": "uuid",
  "email": "colaborador@email.com",
  "role": "colaborador",
  "ativo": "ativo"
}
```

---

## `PUT /authenticated/usuarios`

> Atualiza email e/ou senha do usuário autenticado (roles: `colaborador`, `gerente`, `admin`)

**Request:**
```json
{
  "email": "novo@email.com",
  "senha": "novasenha"
}
```

**Response `200`** (sem body)

---

## `DELETE /authenticated/usuarios`

> Remove o usuário autenticado (roles: `colaborador`, `gerente`, `admin`)

**Response `204`** (sem body)

---

## `GET /admin/all`

> Lista todos os usuários (role: `admin`)

**Response `200`:** array de usuário
