import { apiGet } from "@/lib/api-client";
import type { StarterCodeResponse } from "@/types";
import type { CodeGenLanguage } from "./run-user-solution";

export type { CodeGenLanguage };

export async function getStarterCode(
  problemId: string,
  language: CodeGenLanguage,
  
) {
  return apiGet<StarterCodeResponse>(
    `/${problemId}/starter-code?language=${language}`,
  );
}
