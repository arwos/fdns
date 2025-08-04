-- SEQUENCE
CREATE SEQUENCE IF NOT EXISTS "adblock_link_id_seq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

-- TABLE
CREATE TABLE IF NOT EXISTS "adblock_link" (
	 "id" BIGINT DEFAULT nextval('adblock_link_id_seq') NOT NULL ,
	 CONSTRAINT "adblock_link_id_pk" PRIMARY KEY ("id"),
	 "link" VARCHAR(2049) NOT NULL,
	 "disabled" BOOLEAN NOT NULL,
	 "updated_at" TIMESTAMPTZ NOT NULL,
	 CONSTRAINT "adblock_link__link__uniq" UNIQUE ("link")
);

