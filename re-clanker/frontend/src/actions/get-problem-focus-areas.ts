import { apiGet } from "@/lib/api-client";
import type { ProblemFocusAreasResponse } from "@/types";

export async function getProblemFocusAreas(
  problemId: string,
  
): Promise<ProblemFocusAreasResponse> {
  return apiGet<ProblemFocusAreasResponse>(
    `/${problemId}/focus-areas`,
  );
}
