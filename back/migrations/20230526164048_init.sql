-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
  "id" bigserial primary key,
  "name" text not null,
  "position" text,
  "role" text not null,
  "phone" text,
  "email" text not null unique check(email != ''),
  "password" text not null
);
CREATE TABLE "regions" (
  "id" bigserial primary key,
  "name" text not null unique
);
CREATE TABLE "landfills" (
  "id" bigserial primary key,
  "city" text,
  "name" text,
  "region_id" bigint not null,
  "illegal" bool not null default false,
  "address" text,
  "coordinates" text not null unique check (coordinates != ''),
  "preview_image_path" text,
  "manager" json,
  "type" text,
  "cadastral_number" text,
  "illegal_cadastral_numbers" text,
  "cadastral_category" text,
  "usr_area" text,
  "area" text,
  CONSTRAINT "FK_landfills.region_id" FOREIGN KEY ("region_id") REFERENCES "regions"("id")
);
CREATE TABLE "surveys" (
  "id" bigserial primary key,
  "date" timestamp not null,
  "landfill_id" bigint not null,
  "user_id" bigint,
  CONSTRAINT "FK_surveys.landfill_id" FOREIGN KEY ("landfill_id") REFERENCES "landfills"("id"),
  CONSTRAINT "FK_surveys.user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
CREATE UNIQUE INDEX surveys_date_idx ON public.surveys ("date",landfill_id);
CREATE TABLE "images" (
  "id" bigserial primary key,
  "survey_id" bigint not null,
  "path" text not null,
  "raw_status" bool not null default false,
  CONSTRAINT "FK_images.survey_id" FOREIGN KEY ("survey_id") REFERENCES "surveys"("id")
);
CREATE TABLE "audits" (
  "id" bigserial primary key,
  "date" timestamp not null,
  "user_id" bigint,
  "survey_id" bigint not null,
  "ai_generated_status" bool not null default false,
  CONSTRAINT "FK_audits.survey_id" FOREIGN KEY ("survey_id") REFERENCES "surveys"("id"),
  CONSTRAINT "FK_audits.user_id" FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
CREATE TABLE "violations" (
  "id" bigserial primary key,
  "description" text not null unique,
  "default_status" bool not null default false
);
CREATE TABLE "audits_to_violations" (
  "violation_id" bigint not null,
  "audit_id" bigint not null,
  CONSTRAINT "FK_audits_to_violations.audit_id" FOREIGN KEY ("audit_id") REFERENCES "audits"("id"),
  CONSTRAINT "FK_audits_to_violations.violation_id" FOREIGN KEY ("violation_id") REFERENCES "violations"("id"),
  UNIQUE("violation_id", "audit_id")
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE "audits_to_violations";
DROP TABLE "violations";
DROP TABLE "audits";
DROP TABLE "images";
DROP TABLE "surveys";
DROP TABLE "landfills";
DROP TABLE "regions";
DROP TABLE "session";
DROP TABLE "users";
-- +goose StatementEnd