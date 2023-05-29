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

INSERT INTO accounts (id, balance, currency)
		VALUES ('604f02b2-4e45-48d6-a952-03a0136e8140', 350000, 'EUR');

INSERT INTO accounts (id, balance, currency)
		VALUES ('8fa6c93b-f300-4ef8-9bac-4258caea36db', 500000, 'EUR');
        
INSERT INTO accounts (id, balance, currency)
		VALUES ('ed989ca2-bc1b-413c-8698-d3d9dfa74800', 230000, 'EUR');

INSERT INTO accounts (id, balance, currency)
		VALUES ('6ce82b44-95a5-4e96-915b-1e5b48f3e52a', 120000, 'EUR');

INSERT INTO accounts (id, balance, currency)
		VALUES ('71376d61-8b6c-4289-b5c4-79cb36add23f', 120000, 'EUR');

INSERT INTO transfers (source_account_id, target_account_id, amount, currency)
		VALUES ('8fa6c93b-f300-4ef8-9bac-4258caea36db', '604f02b2-4e45-48d6-a952-03a0136e8140', 50000, 'EUR');
        
INSERT INTO transfers (source_account_id, target_account_id, amount, currency)
		VALUES ('ed989ca2-bc1b-413c-8698-d3d9dfa74800', '6ce82b44-95a5-4e96-915b-1e5b48f3e52a', 70000, 'EUR');