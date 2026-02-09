interface JobDescriptionProps {
  jobDescription: string;
  setJobDescription: (value: string) => void;
}

const JobDescription = ({
  jobDescription,
  setJobDescription,
}: JobDescriptionProps) => {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-1">
        Job Description
      </label>

      <textarea
        rows={6}
        value={jobDescription}
        onChange={(e) => setJobDescription(e.target.value)}
        placeholder="Paste the job description here..."
        className="w-full rounded-xl p-4 text-sm
    bg-white text-gray-900 placeholder-gray-400
    border border-white/30
    focus:outline-none focus:ring-2 focus:ring-emerald-400
    focus:border-emerald-400"
      />
    </div>
  );
};

export default JobDescription;
