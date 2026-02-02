import { apiPost } from "@/lib/api-client";
import type { ProblemText } from "@/types";

export async function generateProblemText(
  problemId: string,
  model: string,
  focusAreaIds?: string[],
) {
  return apiPost<ProblemText>(
    `/api/v1/problems/${problemId}/generate-text`,
    {
      model,
      ...(focusAreaIds && { focusAreaIds }),
    },
  );
}
