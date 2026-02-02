import { apiPost } from "@/lib/api-client";

export async function generateTestCaseOutputs(
  problemId: string,
  model: string,
) {
  return apiPost<{ testCases: unknown[]; jobId: string | null }>(
    `/api/v1/problems/${problemId}/generate-test-case-outputs`,
    { model },
  );
}
