export type GenerationStep =
  | "generateProblemText"
  | "parseFunctionSignature"
  | "generateTestCases"
  | "generateTestCaseInputCode"
  | "generateSolution";

// Step order for sequential execution
export const STEP_ORDER: GenerationStep[] = [
  "generateProblemText",
  "parseFunctionSignature",
  "generateTestCases",
  "generateTestCaseInputCode",
  "generateSolution",
];
