-- SEQUENCE
CREATE SEQUENCE IF NOT EXISTS "adblock_rule_id_seq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

-- TABLE
CREATE TABLE IF NOT EXISTS "adblock_rule" (
	 "id" BIGINT DEFAULT nextval('adblock_rule_id_seq') NOT NULL ,
	 CONSTRAINT "adblock_rule_id_pk" PRIMARY KEY ("id"),
	 "link_id" BIGINT NOT NULL ,
	 CONSTRAINT "adblock_rule_link_id_fk" FOREIGN KEY ("link_id") REFERENCES "adblock_link" ("id") ON DELETE CASCADE NOT DEFERRABLE,
	 "zone" VARCHAR(256) NOT NULL,
	 "rule" VARCHAR(256) NOT NULL,
	 "disabled" BOOLEAN NOT NULL,
	 "updated_at" TIMESTAMPTZ NOT NULL,
	 CONSTRAINT "adblock_rule__rule__uniq" UNIQUE ("rule")
);

