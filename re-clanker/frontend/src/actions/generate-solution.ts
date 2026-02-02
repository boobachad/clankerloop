import { apiPost } from "@/lib/api-client";
import type { Solution } from "@/types";

export async function generateSolution(
  problemId: string,
  model: string,
) {
  return apiPost<Solution>(
    `/api/v1/problems/${problemId}/generate-solution`,
    { model },
  );
}
