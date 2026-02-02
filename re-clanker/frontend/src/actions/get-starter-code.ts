import { apiGet } from "@/lib/api-client";
import type { StarterCodeResponse } from "@/types";

export async function getStarterCode(problemId: string, language: string) {
  return apiGet<StarterCodeResponse>(`/api/v1/problems/${problemId}/starter-code?language=${language}`);
}
