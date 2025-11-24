"use server";

import { Sandbox } from "@/lib/sandbox";
import { getTestCaseInputCode } from "./generate-test-case-input-code";
import { DEFAULT_LANGUAGE } from "@/lib/consts";
import { getProblem, TestCase, updateProblem } from "@/app/api/problem-crud";

export async function generateTestCaseInputs(problemId: string) {
  const testCasesInputCode = await getTestCaseInputCode(problemId);
  const sandbox = await Sandbox.create(DEFAULT_LANGUAGE);

  // Run executions sequentially to avoid conflicts with single-process sandbox
  const results: unknown[] = [];
  for (const testCaseInputCode of testCasesInputCode) {
    // TODO: handle errors
    const result = await sandbox.run(
      testCaseInputCode.inputCode +
        "; const output = generateTestInput();" +
        // Write the output to a file
        "require('fs').writeFileSync('output.json', JSON.stringify(output));"
    );
    console.log("Result of running sandbox code:", result);
    results.push(JSON.parse(await sandbox.readFile("output.json")));
  }

  await sandbox.kill();

  // Merge results into existing test cases
  const { testCases } = await getProblem(problemId);

  const updatedTestCases: TestCase[] = testCases.map(
    (testCase: TestCase, index: number) => {
      const result = results[index];
      if (result === undefined) {
        throw new Error(`Failed to generate result for test case ${index + 1}`);
      }
      if (!Array.isArray(result)) {
        throw new Error(
          `Result for test case ${index + 1} is not an array, got: ${typeof result}`
        );
      }
      return {
        ...testCase,
        input: result,
      };
    }
  );

  // Save updated test cases back to the JSON file
  await updateProblem(problemId, { testCases: updatedTestCases });

  return results;
}

export async function getTestCaseInputs(problemId: string) {
  const { testCases } = await getProblem(problemId);
  return testCases.map((testCase: TestCase) => testCase.input);
}
