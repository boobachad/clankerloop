import { apiPost } from "@/lib/api-client";
import type { TestResult, CustomTestResult, TestCase } from "@/types";

// Re-export types for consumers
export type { TestCase, TestResult, CustomTestResult };

// Define CodeGenLanguage type (shared with hooks)
export type CodeGenLanguage = "typescript" | "python";

export async function runUserSolution(
  problemId: string,
  userCode: string,
  language: CodeGenLanguage = "typescript",
): Promise<TestResult[]> {
  return apiPost<TestResult[]>(
    `/api/v1/problems/${problemId}/solution/run`,
    { code: userCode, language },
  );
}

export async function runUserSolutionWithCustomInputs(
  problemId: string,
  userCode: string,
  customInputs: unknown[][],
  language: CodeGenLanguage = "typescript",
): Promise<CustomTestResult[]> {
  return apiPost<CustomTestResult[]>(
    `/api/v1/problems/${problemId}/solution/run-custom`,
    { code: userCode, customInputs, language },
  );
}
