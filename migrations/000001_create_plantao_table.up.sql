CREATE TABLE plantoes (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    id_colaborador UUID NOT NULL,
    data_inicio TIMESTAMP NOT NULL,
    data_fim TIMESTAMP NOT NULL,
    CONSTRAINT plantoes_pkey PRIMARY KEY(id)
);