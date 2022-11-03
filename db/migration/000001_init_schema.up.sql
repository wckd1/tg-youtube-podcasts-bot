CREATE TABLE "subscriptions" (
  "id" bigserial PRIMARY KEY,
  "channel" varchar NOT NULL,
  "title" varchar NOT NULL,
);