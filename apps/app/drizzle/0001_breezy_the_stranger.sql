ALTER TABLE "user" ADD COLUMN "monthly_points" integer DEFAULT 5 NOT NULL;--> statement-breakpoint
ALTER TABLE "user" ADD COLUMN "lifetime_points" integer DEFAULT 0 NOT NULL;