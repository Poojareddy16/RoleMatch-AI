import type { AnalysisResult } from "../types/analysis";

const API_BASE_URL = "http://localhost:8080";

export const analyzeResume = async (
  resume: File,
  jobDescription: string,
): Promise<AnalysisResult> => {
  const formData = new FormData();
  formData.append("resume", resume);
  formData.append("job_description", jobDescription);

  const response = await fetch(`${API_BASE_URL}/analyze`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    throw new Error("Analysis failed");
  }

  return response.json();
};
