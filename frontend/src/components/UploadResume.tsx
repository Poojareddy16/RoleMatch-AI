interface UploadResumeProps {
  setResume: (file: File | null) => void;
}

const UploadResume = ({ setResume }: UploadResumeProps) => {
  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-1">
        Upload Resume
      </label>

      <div className="flex items-center gap-3">
        <input
          type="file"
          accept=".pdf,.doc,.docx"
          onChange={(e) => setResume(e.target.files?.[0] || null)}
          className="block w-full text-sm text-gray-500
            file:mr-4 file:py-2 file:px-4
            file:rounded-lg file:border-0
            file:text-sm file:font-medium
            file:bg-indigo-50 file:text-indigo-700
            hover:file:bg-indigo-100"
        />
      </div>
    </div>
  );
};

export default UploadResume;
