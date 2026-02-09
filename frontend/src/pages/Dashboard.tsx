import { useState } from "react";
import UploadResume from "../components/UploadResume";
import JobDescription from "../components/JobDescription";
import ResultCard from "../components/ResultCard";
import { analyzeResume } from "../api/api";
import type { AnalysisResult } from "../types/analysis";

const Dashboard = () => {
  const [resume, setResume] = useState<File | null>(null);
  const [jobDescription, setJobDescription] = useState("");
  const [result, setResult] = useState<AnalysisResult | null>(null);
  const [loading, setLoading] = useState(false);

  const handleAnalyze = async () => {
    if (!resume || !jobDescription) return;

    setLoading(true);
    setResult(null);

    try {
      const data = await analyzeResume(resume, jobDescription);
      setResult(data);
    } catch (err) {
      const error = err as Error;
      console.error("Analyze error:", error);
      alert(error.message || "Backend not reachable");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-indigo-600 via-purple-600 to-pink-500 text-white">
      {/* Header */}
      <header className="px-8 py-6">
        <h1 className="text-4xl font-extrabold tracking-tight">RoleMatch-AI</h1>
        <p className="text-indigo-100 mt-1">
          Intelligent Resume ↔ Job Matching
        </p>
      </header>

      {/* Main */}
      <main className="px-8 pb-10 grid grid-cols-1 lg:grid-cols-2 gap-10">
        {/* Input Card */}
        <div className="bg-white/15 backdrop-blur-xl rounded-2xl shadow-2xl p-8 space-y-6 hover:scale-[1.01] transition">
          <h2 className="text-xl font-semibold text-white">Input Details</h2>

          <UploadResume setResume={setResume} />

          <JobDescription
            jobDescription={jobDescription}
            setJobDescription={setJobDescription}
          />

          <button
            onClick={handleAnalyze}
            disabled={!resume || !jobDescription || loading}
            className={`w-full py-3 rounded-xl font-bold tracking-wide transition-all duration-300
              ${
                loading || !resume || !jobDescription
                  ? "bg-white/30 cursor-not-allowed"
                  : "bg-gradient-to-r from-emerald-400 to-cyan-400 text-gray-900 hover:shadow-[0_0_30px_rgba(52,211,153,0.6)] hover:-translate-y-1"
              }`}
          >
            {loading ? "Analyzing with AI…" : "Analyze Resume"}
          </button>
        </div>

        {/* Result Card */}
        <div className="bg-white/15 backdrop-blur-xl rounded-2xl shadow-2xl p-8">
          <h2 className="text-xl font-semibold mb-4">Analysis Result</h2>

          {!result && !loading && (
            <p className="text-indigo-100">
              Upload a resume and job description to see insights.
            </p>
          )}

          {loading && (
            <p className="text-emerald-300 font-semibold animate-pulse">
              AI is thinking…
            </p>
          )}

          {result && <ResultCard result={result} />}
        </div>
      </main>
    </div>
  );
};

export default Dashboard;
