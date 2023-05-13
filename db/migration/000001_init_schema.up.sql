CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "balance" decimal NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "source_account_id" bigint NOT NULL,
  "target_account_id" bigint NOT NULL,
  "amount" decimal NOT NULL,
  "currency" varchar NOT NULL
);

ALTER TABLE "transfers" ADD FOREIGN KEY ("source_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("target_account_id") REFERENCES "accounts" ("id");

CREATE INDEX ON "transfers" ("source_account_id");

CREATE INDEX ON "transfers" ("target_account_id");

CREATE INDEX ON "transfers" ("source_account_id", "target_account_id");

COMMENT ON COLUMN "accounts"."balance" IS 'must always be positive';