-- SEQUENCE
CREATE SEQUENCE IF NOT EXISTS "name_space_id_seq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1;

-- TABLE
CREATE TABLE IF NOT EXISTS "name_space" (
	 "id" BIGINT DEFAULT nextval('name_space_id_seq') NOT NULL ,
	 CONSTRAINT "name_space_id_pk" PRIMARY KEY ("id"),
	 "zone" VARCHAR(256) NOT NULL,
	 "value" TEXT[] NOT NULL,
	 "disabled" BOOLEAN NOT NULL,
	 "updated_at" TIMESTAMPTZ NOT NULL,
	 CONSTRAINT "name_space__zone__uniq" UNIQUE ("zone")
);

