import { apiGet } from "@/lib/api-client";
import type { Model } from "@/types";

export async function listModels() {
  return apiGet<Model[]>("/api/v1/models");
}
