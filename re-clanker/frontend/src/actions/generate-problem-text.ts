import { apiGet, apiPost } from "@/lib/api-client";
import type { ProblemText } from "@/types";

export async function generateProblemText(
  problemId: string,
  model: string,
  enqueueNextStep: boolean = true,
  forceError?: boolean,
  returnDummy?: boolean,
) {
  const data = await apiPost<ProblemText & { jobId: string | null }>(
    `/${problemId}/text/generate`,
    { model, enqueueNextStep, forceError, returnDummy }
  );
  return {
    problemText: data.problemText,
    functionSignature: data.functionSignature,
    problemTextReworded: data.problemTextReworded,
  };
}

export async function getProblemText(
  problemId: string
) {
  return apiGet<ProblemText>(`/${problemId}/text`);
}
