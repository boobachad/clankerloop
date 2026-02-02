import { Suspense } from "react";
import NewProblemPageWrapper from "./problem/[problemId]/components/problem-page-wrapper";

function HomeContent() {
  // No authentication - always show new problem page
  // Note: Focus areas would be fetched client-side via API
  return <NewProblemPageWrapper user={null} focusAreas={[]} />;
}

export default function Home() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <HomeContent />
    </Suspense>
  );
}
