import { apiPost } from "@/lib/api-client";
import type { Model } from "@repo/api-types";

export async function createModel(name: string, encryptedUserId?: string) {
  return apiPost<Model>("/models", { name }, encryptedUserId);
}
