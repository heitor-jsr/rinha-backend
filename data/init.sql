CREATE SEQUENCE public.client_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.client_id_seq OWNER TO postgres;

SET default_tablespace = '';
SET default_table_access_method = heap;

CREATE TABLE public.clientes (
    id integer DEFAULT nextval('public.client_id_seq'::regclass) NOT NULL,
    limite integer,
    saldo integer,
    nome varchar(255),
    PRIMARY KEY (id)
);

ALTER TABLE public.clientes OWNER TO postgres;

CREATE TABLE public.saldo (
    total integer,
    data_extrato timestamp without time zone DEFAULT NOW(),
    limite integer
);

ALTER TABLE public.saldo OWNER TO postgres;

CREATE TABLE public.transacoes (
    id integer DEFAULT nextval('public.client_id_seq'::regclass) NOT NULL,
    PRIMARY KEY (id),
    valor integer,
    tipo varchar(255),
    descricao varchar(255),
    realizada_em timestamp without time zone DEFAULT NOW(),
	cliente_id INTEGER NOT NULL,
	CONSTRAINT fk_clientes_transacoes_id
		FOREIGN KEY (cliente_id) REFERENCES clientes(id)
);

ALTER TABLE public.transacoes OWNER TO postgres;

INSERT INTO public.clientes (limite, saldo) VALUES
(100000, 0),
(80000, 0),
(1000000, 0),
(10000000, 0);
