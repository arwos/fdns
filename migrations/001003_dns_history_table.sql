-- SEQUENCE
CREATE SEQUENCE IF NOT EXISTS "dns_history_id_seq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

-- TABLE
CREATE TABLE IF NOT EXISTS "dns_history" (
	 "id" BIGINT DEFAULT nextval('dns_history_id_seq') NOT NULL ,
	 CONSTRAINT "dns_history_id_pk" PRIMARY KEY ("id"),
	 "domain" VARCHAR(256) NOT NULL,
	 "qtype" SMALLINT NOT NULL,
	 "value" TEXT[] NOT NULL,
	 "disabled" BOOLEAN NOT NULL,
	 "updated_at" TIMESTAMPTZ NOT NULL,
	 CONSTRAINT "dns_history__domain_qtype__uniq" UNIQUE ("domain","qtype")
);

