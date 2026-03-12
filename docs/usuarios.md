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

## `GET /authenticated/usuarios/:id_usuario`
> Autenticado (role: `colaborador`)

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

## `PUT /authenticated/usuarios/:id_usuario`
> Autenticado (role: `colaborador`)

**Request:**
```json
{
  "email": "novo@email.com",
  "senha": "novasenha"
}
```

**Response `200`** (sem body)

---

## `DELETE /authenticated/usuarios/:id_usuario`
> Autenticado (role: `colaborador`)

**Response `204`** (sem body)