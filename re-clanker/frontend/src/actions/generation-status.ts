import { apiGet } from "@/lib/api-client";
import type { GenerationStatus } from "@/types";

// Re-export types for consumers
export type { GenerationStep, GenerationStatus } from "@/types";

export async function getGenerationStatus(
  problemId: string,
  
): Promise<GenerationStatus> {
  return apiGet<GenerationStatus>(
    `/${problemId}/generation-status`,
  );
}
