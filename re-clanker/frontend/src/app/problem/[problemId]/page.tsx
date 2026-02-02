"use server";
import { redirect } from "next/navigation";
import ProblemRender from "./components/problem-render";

export default async function Page({
  params,
}: {
  params: Promise<{ problemId: string }>;
}) {
  const { problemId } = await params;
  
  // No auth - render problem directly
  return <ProblemRender problemId={problemId} user={null} />;
}
