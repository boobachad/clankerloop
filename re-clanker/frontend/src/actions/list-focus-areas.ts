import { apiGet } from "@/lib/api-client";
import type { FocusArea } from "@/types";

export async function listFocusAreas() {
  return apiGet<FocusArea[]>("/api/v1/focus-areas");
}
