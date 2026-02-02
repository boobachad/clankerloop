import { apiGet, apiPost } from "@/lib/api-client";

interface TestOutputsGenerateResponse {
  testCases: unknown[];
  jobId: string | null;
}

export async function generateTestCaseOutputs(
  problemId: string,
  
  enqueueNextStep: boolean = true,
) {
  const data = await apiPost<TestOutputsGenerateResponse>(
    `/${problemId}/test-cases/outputs/generate`,
    { enqueueNextStep },
  );
  return data.testCases;
}

export async function getTestCaseOutputs(
  problemId: string,
  
) {
  return apiGet<unknown[]>(`/${problemId}/test-cases/outputs`);
}
