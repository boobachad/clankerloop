import { apiGet } from "@/lib/api-client";
import type { StarterCodeResponse } from "@repo/api-types";

export type CodeGenLanguage = "typescript" | "python";

export async function getStarterCode(
  problemId: string,
  language: CodeGenLanguage,
  encryptedUserId?: string,
) {
  return apiGet<StarterCodeResponse>(
    `/${problemId}/starter-code?language=${language}`,
    encryptedUserId,
  );
}
