import { apiGet, apiPost } from "@/lib/api-client";
import type { TestCase, TestCaseDescription } from "@/types";

interface TestCasesGenerateResponse {
  testCases: TestCaseDescription[];
  jobId: string | null;
}

export async function generateTestCases(
  problemId: string,
  model: string,
  
  enqueueNextStep: boolean = true,
  forceError?: boolean,
  returnDummy?: boolean,
) {
  const data = await apiPost<TestCasesGenerateResponse>(
    `/${problemId}/test-cases/generate`,
    { model, enqueueNextStep, forceError, returnDummy },
  );
  return data.testCases;
}

export async function getTestCases(
  problemId: string,
  
) {
  return apiGet<TestCase[]>(`/${problemId}/test-cases`);
}
