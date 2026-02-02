import { apiPost } from "@/lib/api-client";

export async function generateTestCaseInputCode(
  problemId: string,
  model: string,
) {
  return apiPost<{ inputCodes: string[]; jobId: string | null }>(
    `/api/v1/problems/${problemId}/generate-test-case-input-code`,
    { model },
  );
}
