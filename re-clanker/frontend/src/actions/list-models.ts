import { apiGet } from "@/lib/api-client";
import type { Model } from "@/types";

export async function listModels(encryptedUserId?: string) {
  return apiGet<Model[]>("/models");
}
