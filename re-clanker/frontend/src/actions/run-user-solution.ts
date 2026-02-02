import { apiPost } from "@/lib/api-client";
import type { TestResult, CustomTestResult } from "@/types";

// Re-export types for consumers
export type { TestCase, TestResult, CustomTestResult } from "@/types";

// Define CodeGenLanguage type (shared with hooks)
export type CodeGenLanguage = "typescript" | "python";

export async function runUserSolution(
  problemId: string,
  userCode: string,
  language: CodeGenLanguage = "typescript",
  
): Promise<TestResult[]> {
  return apiPost<TestResult[]>(
    `/${problemId}/solution/run`,
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
    `/${problemId}/solution/run-custom`,
    { code: userCode, customInputs, language },
  );
}
