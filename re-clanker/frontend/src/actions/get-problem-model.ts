import { apiGet } from "@/lib/api-client";

export async function getProblemModel(
  problemId: string,
  
) {
  const data = await apiGet<{ model: string | null }>(
    `/${problemId}/model`,
  );
  return data.model;
}
