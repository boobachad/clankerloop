import { backendGet, backendPost } from "@/lib/backend-client";

export async function generateSolution(
  problemId: string,
  model: string,
  encryptedUserId?: string,
  updateProblem: boolean = true,
  enqueueNextStep: boolean = true
) {
  const response = await backendPost<{ solution: string; jobId: string }>(
    `/problems/${problemId}/solution/generate`,
    { model, updateProblem, enqueueNextStep },
    encryptedUserId
  );
  return response.solution;
}

export async function getSolution(problemId: string, encryptedUserId?: string) {
  return backendGet<string>(`/problems/${problemId}/solution`, encryptedUserId);
}
