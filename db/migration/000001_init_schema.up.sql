CREATE TABLE "accounts" (
  "id" uuid DEFAULT gen_random_uuid(),
  "balance" decimal NOT NULL,
  "currency" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  PRIMARY KEY ("id")  
);

CREATE TABLE "transfers" (
  "id" uuid DEFAULT gen_random_uuid(),
  "source_account_id" uuid NOT NULL,
  "target_account_id" uuid NOT NULL,
  "amount" decimal NOT NULL,
  "currency" varchar NOT NULL,
  PRIMARY KEY ("id")  
);

ALTER TABLE "transfers" ADD FOREIGN KEY ("source_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("target_account_id") REFERENCES "accounts" ("id");