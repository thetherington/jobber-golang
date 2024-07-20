CREATE TABLE "auth" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar NOT NULL CONSTRAINT username_unique UNIQUE,
  "password" varchar NOT NULL,
  "profile_public_id" varchar NOT NULL,
  "email" varchar NOT NULL CONSTRAINT email_unique UNIQUE,
  "country" varchar NOT NULL,
  "profile_picture" varchar NOT NULL,
  "email_verification_token" varchar NOT NULL,
  "email_verified" boolean NOT NULL DEFAULT FALSE,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "password_reset_token" varchar,
  "password_reset_expires" timestamptz NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX auth_unique ON "auth" ("username", "email", "email_verification_token");