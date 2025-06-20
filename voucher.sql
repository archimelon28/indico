-- Table: public.voucher

-- DROP TABLE IF EXISTS public.voucher;

CREATE TABLE IF NOT EXISTS public.voucher
(
    id integer NOT NULL DEFAULT nextval('voucher_id_seq'::regclass),
    voucher_code character varying(20) COLLATE pg_catalog."default" NOT NULL,
    discount_percent integer NOT NULL,
    expiry_date date NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    status integer,
    CONSTRAINT voucher_pkey PRIMARY KEY (id)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.voucher
    OWNER to postgres;