import { apiPost } from "@/lib/api-client";
import type { CreateProblemResponse, StartFrom } from "@/types";

export async function createProblem(
  model: string,
  autoGenerate: boolean = true,
  returnDummy?: boolean,
  startFrom?: StartFrom,
  focusAreaIds?: string[],
) {
  return apiPost<CreateProblemResponse>(
    `/api/v1/problems`,
    {
      model,
      ...(returnDummy !== undefined && { returnDummy }),
      ...(startFrom !== undefined && { startFrom }),
      ...(focusAreaIds !== undefined &&
        focusAreaIds.length > 0 && { focusAreaIds }),
    },
  );
}
