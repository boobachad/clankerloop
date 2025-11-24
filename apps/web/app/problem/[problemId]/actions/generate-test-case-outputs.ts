"use server";

import { Sandbox } from "@/lib/sandbox";
import { DEFAULT_LANGUAGE } from "@/lib/consts";
import { getSolution } from "./generate-solution";
import { getTestCases } from "./generate-test-cases";
import { getProblem, TestCase, updateProblem } from "@/app/api/problem-crud";

export async function generateTestCaseOutputs(problemId: string) {
  const solution = await getSolution(problemId);
  const testCaseInputs = await getTestCases(problemId);
  const sandbox = await Sandbox.create(DEFAULT_LANGUAGE);

  // Run executions sequentially to avoid conflicts with single-process sandbox
  const results: unknown[] = [];
  for (const testCaseInput of testCaseInputs) {
    // TODO: handle errors
    await sandbox.run(
      solution +
        "; const output = runSolution(..." +
        JSON.stringify(testCaseInput.input) +
        ");" +
        // Write the output to a file
        "require('fs').writeFileSync('output.json', JSON.stringify(output));"
    );
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
      return {
        ...testCase,
        expected: result,
      };
    }
  );

  // Save updated test cases back to the JSON file
  await updateProblem(problemId, { testCases: updatedTestCases });

  return results;
}

export async function getTestCaseOutputs(problemId: string) {
  const { testCases } = await getProblem(problemId);
  return testCases.map((testCase: TestCase) => testCase.expected);
}
