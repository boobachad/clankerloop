import { apiGet } from "@/lib/api-client";
import type { ProblemModel } from "@/types";

export async function getProblemModel(problemId: string) {
  return apiGet<ProblemModel>(`/api/v1/problems/${problemId}/model`);
}
