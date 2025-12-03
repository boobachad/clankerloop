ALTER TABLE "problems" ADD COLUMN "easier_than" uuid;--> statement-breakpoint
ALTER TABLE "problems" ADD COLUMN "harder_than" uuid;--> statement-breakpoint
ALTER TABLE "problems" ADD CONSTRAINT "problems_easier_than_fkey" FOREIGN KEY ("easier_than") REFERENCES "problems"("id") ON DELETE SET NULL;--> statement-breakpoint
ALTER TABLE "problems" ADD CONSTRAINT "problems_harder_than_fkey" FOREIGN KEY ("harder_than") REFERENCES "problems"("id") ON DELETE SET NULL;