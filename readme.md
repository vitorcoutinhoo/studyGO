## Documentação da API 

### Passo 1:

Criar o Banco Postegree com nome `Plantão`, criar as tabelas das entidades no banco.

### Passo 2:

Se for rodar local utilizar os comandos no terminal: 
```bash
go mod tidy
go run main.go
``` 
Caso for o docker compose rodar o seguinte comando no terminal: 
```bash
docker compose up --build
```

### Passo 3:

Acesse o link `http://localhost:8080/` a partir dele, terá acesso aos endpoints.