import { apiGet, apiPost } from "@/lib/api-client";
import type { Solution } from "@/types";

interface SolutionGenerateResponse {
  solution: string | null;
  jobId: string | null;
}

export async function generateSolution(
  problemId: string,
  model: string,
  
  updateProblem: boolean = true,
  enqueueNextStep: boolean = true,
  forceError?: boolean,
  returnDummy?: boolean,
) {
  const data = await apiPost<SolutionGenerateResponse>(
    `/${problemId}/solution/generate`,
    { model, updateProblem, enqueueNextStep, forceError, returnDummy },
  );
  return data.solution;
}

export async function getSolution(problemId: string) {
  const data = await apiGet<Solution>(
    `/${problemId}/solution`,
  );
  return data.solution;
}
