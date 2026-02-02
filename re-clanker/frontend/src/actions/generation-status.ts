import { apiGet } from "@/lib/api-client";
import type { GenerationStatus } from "@/types";

export async function getGenerationStatus(problemId: string) {
  return apiGet<GenerationStatus>(`/api/v1/problems/${problemId}/generation-status`);
}
