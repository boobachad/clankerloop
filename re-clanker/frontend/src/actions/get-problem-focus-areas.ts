import { apiGet } from "@/lib/api-client";
import type { ProblemFocusAreasResponse } from "@/types";

export async function getProblemFocusAreas(problemId: string) {
  return apiGet<ProblemFocusAreasResponse>(`/api/v1/problems/${problemId}/focus-areas`);
}
