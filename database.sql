CREATE TABLE colaboradores (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    telefone VARCHAR(50),
    cargo VARCHAR(100),
    departamento VARCHAR(100),
    foto_url TEXT,
    ativo CHAR(1) DEFAULT 'Y',
    data_admissao DATE,
    data_desligamento DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT colaboradores_pkey PRIMARY KEY (id)
);

CREATE TABLE config_valores_dia (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    tipo_dia VARCHAR(100) NOT NULL UNIQUE,
    valor NUMERIC NOT NULL,
    descricao TEXT,
    vigencia_inicio DATE NOT NULL,
    vigencia_fim DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT config_valores_dia_pkey PRIMARY KEY (id)
);

CREATE TABLE modelos_comunicacao (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    nome VARCHAR(255) NOT NULL UNIQUE,
    tipo VARCHAR(100) NOT NULL,
    assunto VARCHAR(255),
    corpo TEXT NOT NULL,
    ativo CHAR(1) DEFAULT 'Y',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT modelos_comunicacao_pkey PRIMARY KEY (id)
);

CREATE TABLE feriados (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    data DATE NOT NULL UNIQUE,
    nome VARCHAR(255) NOT NULL,
    descricao TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT feriados_pkey PRIMARY KEY (id)
);

CREATE TABLE plantoes (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    id_colaborador UUID NOT NULL,
    data_inicio TIMESTAMP NOT NULL,
    data_fim TIMESTAMP NOT NULL,
    status VARCHAR(50) DEFAULT 'agendado',
    valor_total NUMERIC DEFAULT 0,
    observacoes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT plantoes_pkey PRIMARY KEY (id),
    CONSTRAINT plantoes_id_colaborador_fkey
        FOREIGN KEY (id_colaborador)
        REFERENCES colaboradores (id)
);

CREATE TABLE plantoes_detalhes (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    id_plantao UUID NOT NULL,
    data DATE NOT NULL,
    tipo_dia VARCHAR(100) NOT NULL,
    valor NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT plantoes_detalhes_pkey PRIMARY KEY (id),
    CONSTRAINT plantoes_detalhes_id_plantao_fkey
        FOREIGN KEY (id_plantao)
        REFERENCES plantoes (id)
);

CREATE TABLE pagamentos (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    id_plantao UUID NOT NULL,
    id_colaborador UUID NOT NULL,
    valor_total NUMERIC(10,2) NOT NULL,
    data_pagamento DATE,
    status VARCHAR(50) DEFAULT 'pendente',
    observacoes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT pagamentos_pkey PRIMARY KEY (id),
    CONSTRAINT pagamentos_id_plantao_fkey
        FOREIGN KEY (id_plantao)
        REFERENCES plantoes (id),
    CONSTRAINT pagamentos_id_colaborador_fkey
        FOREIGN KEY (id_colaborador)
        REFERENCES colaboradores (id)
);

CREATE TABLE status_plantao (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    id_plantao UUID NOT NULL,
    status_antigo VARCHAR(50),
    status_novo VARCHAR(50) NOT NULL,
    id_usuario UUID,
    observacoes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT status_plantao_pkey PRIMARY KEY (id),
    CONSTRAINT status_plantao_id_plantao_fkey
        FOREIGN KEY (id_plantao)
        REFERENCES plantoes (id)
);

CREATE TABLE envios_comunicacao (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    id_modelo UUID,
    id_colaborador UUID NOT NULL,
    tipo VARCHAR(100) NOT NULL,
    destinatario VARCHAR(255) NOT NULL,
    assunto VARCHAR(255),
    corpo TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'enviado',
    data_envio TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_leitura TIMESTAMP,
    erro_log TEXT,
    CONSTRAINT envios_comunicacao_pkey PRIMARY KEY (id),
    CONSTRAINT envios_comunicacao_id_modelo_fkey
        FOREIGN KEY (id_modelo)
        REFERENCES modelos_comunicacao (id),
    CONSTRAINT envios_comunicacao_id_colaborador_fkey
        FOREIGN KEY (id_colaborador)
        REFERENCES colaboradores (id)
);

CREATE TABLE usuarios_login (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    id_colaborador UUID NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    senha_hash TEXT NOT NULL,
    role VARCHAR(50) DEFAULT 'colaborador',
    ativo CHAR(1) DEFAULT 'Y',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT usuarios_login_pkey PRIMARY KEY (id),
    CONSTRAINT usuarios_login_id_colaborador_fkey
        FOREIGN KEY (id_colaborador)
        REFERENCES colaboradores (id)
);
