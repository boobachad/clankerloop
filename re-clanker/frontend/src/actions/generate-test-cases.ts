import { apiPost } from "@/lib/api-client";
import type { TestCaseDescription } from "@/types";

export async function generateTestCases(
  problemId: string,
  model: string,
) {
  return apiPost<{ testCases: TestCaseDescription[]; jobId: string | null }>(
    `/api/v1/problems/${problemId}/generate-test-cases`,
    { model },
  );
}
