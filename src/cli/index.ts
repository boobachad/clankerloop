import type { Difficulty, Language } from "../types/index.js";
import { generateCommand, type GenerateOptions } from "./generate.js";
import { solveCommand, type SolveOptions } from "./solve.js";

export { generateCommand, solveCommand };
export type { GenerateOptions, SolveOptions };

/**
 * Parse command line arguments and route to appropriate command
 */
export async function runCLI(args: string[]): Promise<void> {
  const command = args[0];

  if (command === "generate") {
    await handleGenerateCommand(args.slice(1));
  } else if (command === "solve") {
    await handleSolveCommand(args.slice(1));
  } else {
    printHelp();
    process.exit(1);
  }
}

function handleGenerateCommand(args: string[]): Promise<void> {
  const options: GenerateOptions = {
    model: getArg(args, "--model") || "google/gemini-2.0-flash",
    difficulty: (getArg(args, "--difficulty") as Difficulty) || "medium",
    language: (getArg(args, "--language") as Language) || "javascript",
    topic: getArg(args, "--topic"),
    numTestCases: parseInt(getArg(args, "--tests") || "10"),
    numSamples: parseInt(getArg(args, "--samples") || "3"),
    output: getArg(args, "--output"),
  };

  return generateCommand(options);
}

function handleSolveCommand(args: string[]): Promise<void> {
  const problemFile = getArg(args, "--problem");
  const solutionFile = getArg(args, "--solution");

  if (!problemFile || !solutionFile) {
    console.error(
      "Error: --problem and --solution are required for solve command"
    );
    printHelp();
    process.exit(1);
  }

  const options: SolveOptions = {
    problemFile,
    solutionFile,
    language: (getArg(args, "--language") as Language) || "javascript",
    showHidden: args.includes("--show-hidden"),
  };

  return solveCommand(options);
}

function getArg(args: string[], flag: string): string | undefined {
  const index = args.indexOf(flag);
  if (index === -1 || index === args.length - 1) return undefined;
  return args[index + 1];
}

function printHelp(): void {
  console.log(`
AI LeetCode Generator - Usage:

GENERATE A PROBLEM:
  bun run index.ts generate [options]

  Options:
    --model <string>        AI model to use (default: "google/gemini-2.0-flash")
    --difficulty <level>    Problem difficulty: easy, medium, hard (default: medium)
    --language <lang>       Target language: javascript, typescript, python (default: javascript)
    --topic <string>        Problem topic (e.g., "arrays", "dynamic programming")
    --tests <number>        Number of test cases to generate (default: 10)
    --samples <number>      Number of sample test cases (default: 3)
    --output <file>         Save problem to JSON file

  Example:
    bun run index.ts generate --model "google/gemini-2.0-flash" --difficulty medium --output problem.json

TEST A SOLUTION:
  bun run index.ts solve [options]

  Options:
    --problem <file>        Path to problem JSON file (required)
    --solution <file>       Path to solution file (required)
    --language <lang>       Solution language: javascript, typescript, python (required)
    --show-hidden           Show results of hidden test cases

  Example:
    bun run index.ts solve --problem problem.json --solution solution.js --language javascript
`);
}
