import { apiPost } from "@/lib/api-client";
import type { CreateProblemResponse } from "@repo/api-types";

export async function createProblem(model: string, encryptedUserId?: string) {
  return apiPost<CreateProblemResponse>("", { model }, encryptedUserId);
}
