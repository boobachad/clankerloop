import { apiGet } from "@/lib/api-client";
import type { FocusArea } from "@/types";

export async function listFocusAreas(
  
): Promise<FocusArea[]> {
  return apiGet<FocusArea[]>("/focus-areas");
}
