import * as React from "react";

interface Props {
  message?: string;
}

const Splash: React.FC<Props> = ({ message = "Loading..." }) => {
  return (
    <div className="fixed inset-0 flex flex-col items-center justify-center bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100 z-50">
      {/* Spinner */}
      <div
        className="animate-spin h-16 w-16 mb-4 rounded-full 
                      border-t-4 border-b-4 border-zinc-300 border-t-zinc-600 
                      shadow-[0_2px_4px_rgba(0,0,0,0.2)]"
      ></div>

      {/* Message */}
      <p className="text-zinc-700 text-lg font-semibold select-none">
        {message}
      </p>
    </div>
  );
};

export default Splash;
