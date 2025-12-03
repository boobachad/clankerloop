CREATE TYPE "public"."user_problem_attempt_status" AS ENUM('attempt', 'run', 'pass');--> statement-breakpoint
CREATE TABLE "user_problem_attempts" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"user_id" text NOT NULL,
	"problem_id" uuid NOT NULL,
	"submission_code" text NOT NULL,
	"submission_language" text NOT NULL,
	"status" "user_problem_attempt_status" DEFAULT 'attempt' NOT NULL,
	"created_at" timestamp DEFAULT now() NOT NULL,
	"updated_at" timestamp DEFAULT now() NOT NULL
);
--> statement-breakpoint
ALTER TABLE "user_problem_attempts" ADD CONSTRAINT "user_problem_attempts_problem_id_problems_id_fk" FOREIGN KEY ("problem_id") REFERENCES "public"."problems"("id") ON DELETE cascade ON UPDATE no action;