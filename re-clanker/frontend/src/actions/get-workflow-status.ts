import { apiGet } from "@/lib/api-client";
import type { WorkflowStatusResponse } from "@/types";

export async function getWorkflowStatus(problemId: string) {
  return apiGet<WorkflowStatusResponse>(`/api/v1/problems/${problemId}/workflow-status`);
}
