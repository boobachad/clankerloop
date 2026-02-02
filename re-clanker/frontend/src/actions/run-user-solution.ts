import { apiPost } from "@/lib/api-client";
import type { TestResult, CustomTestResult } from "@/types";

export async function runUserSolution(
  problemId: string,
  code: string,
  language: string,
  customTests?: unknown[][],
) {
  if (customTests) {
    // Run with custom tests
    return apiPost<CustomTestResult[]>(
      `/api/v1/problems/${problemId}/run-custom-tests`,
      {
        code,
        language,
        customInputs: customTests,
      },
    );
  } else {
    // Run with problem test cases
    return apiPost<TestResult[]>(
      `/api/v1/problems/${problemId}/run-solution`,
      {
        code,
        language,
      },
    );
  }
}
