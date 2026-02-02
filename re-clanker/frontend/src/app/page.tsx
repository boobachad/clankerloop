"use server";
import NewProblemPageWrapper from "./problem/[problemId]/components/problem-page-wrapper";

export default async function Home() {
  // No authentication - always show new problem page
  // Note: Focus areas would be fetched client-side via API
  return <NewProblemPageWrapper user={null} focusAreas={[]} />;
}
