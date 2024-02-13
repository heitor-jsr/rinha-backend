-- Criação do usuário
CREATE USER postgres WITH PASSWORD 'password';

-- Concede permissões ao usuário
GRANT ALL PRIVILEGES ON DATABASE rinha TO postgres;

-- Sequência para gerar IDs automaticamente para a tabela de clientes
CREATE SEQUENCE public.client_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.client_id_seq OWNER TO postgres;

SET default_tablespace = '';
SET default_table_access_method = heap;

-- Tabela de clientes
CREATE TABLE public.clients (
    id integer DEFAULT nextval('public.client_id_seq'::regclass) NOT NULL,
    limite integer,
    saldo integer,
    PRIMARY KEY (id)
);

ALTER TABLE public.clients OWNER TO postgres;

-- Tabela de extratos
CREATE TABLE public.statements (
    id SERIAL PRIMARY KEY,
    total integer,
    statement_date timestamp without time zone,
    limite integer
);

ALTER TABLE public.statements OWNER TO postgres;

-- Tabela de transações
CREATE TABLE public.transactions (
    valor integer,
    tipo varchar(255),
    descricao varchar(255),
    done_at timestamp without time zone,
    client_id integer REFERENCES public.clients(id)
);

ALTER TABLE public.transactions OWNER TO postgres;

-- Inserção de dados na tabela de clientes
INSERT INTO public.clients (limite, saldo) VALUES
(100000, 0),
(80000, 0),
(1000000, 0),
(10000000, 0),
(500000, 0);
