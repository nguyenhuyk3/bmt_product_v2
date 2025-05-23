CREATE TYPE "genres" AS ENUM (
  'action',
  'adventure',
  'animation',
  'comedy',
  'crime',
  'drama',
  'fantasy',
  'historical',
  'horror',
  'mystery',
  'romance',
  'sci_fi',
  'thriller',
  'war',
  'western',
  'dark_comedy',
  'documentary',
  'musical',
  'sports',
  'superhero',
  'psychological_thriller',
  'slasher',
  'biopic',
  'noir',
  'family'
);

CREATE TYPE "languages" AS ENUM (
  'vietnamese',
  'english'
);

CREATE TYPE "statuses" AS ENUM (
  'failed',
  'pending',
  'success'
);

CREATE TYPE "fab_types" AS ENUM (
  'food',
  'beverage'
);

CREATE TABLE "films" (
  "id" serial PRIMARY KEY NOT NULL,
  "title" text NOT NULL,
  "description" text NOT NULL,
  "release_date" date NOT NULL,
  "duration" interval NOT NULL
);

CREATE TABLE "fillm_changes" (
  "film_id" int NOT NULL,
  "changed_by" varchar(32) NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "film_genres" (
  "film_id" int,
  "genre" genres
);

CREATE TABLE "other_film_informations" (
  "film_id" int PRIMARY KEY,
  "status" statuses,
  "poster_url" text,
  "trailer_url" text
);

CREATE TABLE "foods_and_beverages" (
  "id" serial PRIMARY KEY NOT NULL,
  "name" text NOT NULL,
  "type" fab_types NOT NULL,
  "image_url" text,
  "price" int NOT NULL DEFAULT 0,
  "is_deleted" boolean NOT NULL DEFAULT false,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now())
);

CREATE TABLE "outboxes" (
  "id" uuid PRIMARY KEY NOT NULL DEFAULT (gen_random_uuid()),
  "aggregated_type" varchar(64) NOT NULL,
  "aggregated_id" int NOT NULL,
  "event_type" varchar(64) NOT NULL,
  "payload" jsonb NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "films" ("id");

CREATE UNIQUE INDEX ON "film_genres" ("film_id", "genre");

CREATE UNIQUE INDEX ON "other_film_informations" ("film_id");

CREATE INDEX ON "foods_and_beverages" ("id");

CREATE INDEX ON "foods_and_beverages" ("name");

CREATE INDEX ON "outboxes" ("aggregated_type", "aggregated_id");

ALTER TABLE "fillm_changes" ADD FOREIGN KEY ("film_id") REFERENCES "films" ("id");

ALTER TABLE "film_genres" ADD FOREIGN KEY ("film_id") REFERENCES "films" ("id");

ALTER TABLE "other_film_informations" ADD FOREIGN KEY ("film_id") REFERENCES "films" ("id");

CREATE PUBLICATION product_dbz_publication FOR TABLE outboxes;