import { redirect } from "next/navigation";
import CreateProblemForm from "./components/create-problem-form";

export default async function CreateProblemPage() {
  // No auth - allow anyone to create problems
  return <CreateProblemForm />;
}
