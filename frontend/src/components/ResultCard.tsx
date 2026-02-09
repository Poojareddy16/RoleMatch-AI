import type { AnalysisResult } from "../types/analysis";

interface Props {
  result: AnalysisResult;
}

const ResultCard = ({ result }: Props) => {
  const scoreColor =
    result.match_score >= 75
      ? "from-emerald-400 to-lime-400"
      : result.match_score >= 50
        ? "from-amber-400 to-orange-400"
        : "from-red-400 to-pink-400";

  return (
    <div className="space-y-6">
      {/* Match Score */}
      <div>
        <p className="text-sm text-indigo-100 mb-1">Match Score</p>

        <div className="flex items-center gap-4">
          <div
            className={`text-4xl font-extrabold bg-gradient-to-r ${scoreColor} bg-clip-text text-transparent`}
          >
            {result.match_score}%
          </div>

          <div className="flex-1 h-3 rounded-full bg-white/20 overflow-hidden">
            <div
              className={`h-3 rounded-full bg-gradient-to-r ${scoreColor} transition-all duration-700`}
              style={{ width: `${result.match_score}%` }}
            />
          </div>
        </div>
      </div>

      {/* Missing Skills */}
      <div>
        <p className="text-sm text-indigo-100 mb-2">Missing Skills</p>

        <div className="flex flex-wrap gap-2">
          {result.missing_skills.map((skill, idx) => (
            <span
              key={idx}
              className="px-3 py-1 text-xs font-semibold rounded-full
                bg-red-500/20 text-red-200
                hover:bg-red-500 hover:text-white
                transition"
            >
              {skill}
            </span>
          ))}
        </div>
      </div>

      {/* Summary */}
      <div>
        <p className="text-sm text-indigo-100 mb-1">AI Summary</p>
        <p className="text-sm leading-relaxed text-white/90">
          {result.summary}
        </p>
      </div>
    </div>
  );
};

export default ResultCard;
