-- SEQUENCE
CREATE SEQUENCE IF NOT EXISTS "static_rule_id_seq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

-- TABLE
CREATE TABLE IF NOT EXISTS "static_rule" (
	 "id" BIGINT DEFAULT nextval('static_rule_id_seq') NOT NULL ,
	 CONSTRAINT "static_rule_id_pk" PRIMARY KEY ("id"),
	 "rule" VARCHAR(256) NOT NULL,
	 "qtype" SMALLINT NOT NULL,
	 "value" TEXT[] NOT NULL,
	 "disabled" BOOLEAN NOT NULL,
	 "updated_at" TIMESTAMPTZ NOT NULL,
	 CONSTRAINT "static_rule__rule_qtype__uniq" UNIQUE ("rule","qtype")
);

