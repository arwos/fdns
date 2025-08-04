-- SEQUENCE
CREATE SEQUENCE IF NOT EXISTS "api_token_id_seq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

-- TABLE
CREATE TABLE IF NOT EXISTS "api_token" (
	 "id" BIGINT DEFAULT nextval('api_token_id_seq') NOT NULL ,
	 CONSTRAINT "api_token_id_pk" PRIMARY KEY ("id"),
	 "token" UUID NOT NULL,
	 "qtype" SMALLINT NOT NULL,
	 "domain" TEXT NOT NULL,
	 "disabled" BOOLEAN NOT NULL,
	 "updated_at" TIMESTAMPTZ NOT NULL,
	 CONSTRAINT "api_token__token_domain__uniq" UNIQUE ("token","domain")
);

