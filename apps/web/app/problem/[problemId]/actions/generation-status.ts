import { backendGet } from "@/lib/backend-client";

export type GenerationStep =
  | "generateProblemText"
  | "generateTestCases"
  | "generateTestCaseInputCode"
  | "generateTestCaseInputs"
  | "generateSolution"
  | "generateTestCaseOutputs";

export interface GenerationStatus {
  jobId?: string;
  status: "none" | "pending" | "in_progress" | "completed" | "failed";
  currentStep?: GenerationStep;
  completedSteps?: GenerationStep[];
  progress?: {
    completed: number;
    total: number;
    percent: number;
  };
  error?: string;
}

export async function getGenerationStatus(
  problemId: string,
  encryptedUserId?: string
): Promise<GenerationStatus> {
  return backendGet<GenerationStatus>(
    `/problems/${problemId}/generation-status`,
    encryptedUserId
  );
}
