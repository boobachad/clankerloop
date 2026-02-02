import { apiGet } from "@/lib/api-client";
import type { WorkflowStatusResponse } from "@/types";

// Re-export types for consumers
export type { WorkflowStatus, WorkflowStatusResponse } from "@/types";

export async function getWorkflowStatus(
  problemId: string,
  
): Promise<WorkflowStatusResponse> {
  return apiGet<WorkflowStatusResponse>(
    `/${problemId}/workflow-status`,
  );
}
