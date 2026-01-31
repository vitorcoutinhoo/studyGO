CREATE TABLE plantoes (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    colaborador_id UUID NOT NULL,
    inicio TIMESTAMP NOT NULL,
    fim TIMESTAMP NOT NULL,
    status integer not null,
    CONSTRAINT plantoes_pkey PRIMARY KEY(id)
);