import * as React from "react";
import { Link } from "react-router-dom";

const NotFound: React.FC = () => {
  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100">
      {/* 404 number */}
      <h1 className="text-6xl font-bold mb-4 text-zinc-700 drop-shadow-[0_2px_4px_rgba(0,0,0,0.2)] select-none">
        404
      </h1>

      {/* Message */}
      <p className="text-lg mb-6 text-zinc-600 select-none">Page not found</p>

      {/* Back home link styled as button */}
      <Link
        to="/"
        className="h-10 px-5 rounded bg-gradient-to-b from-zinc-100 to-zinc-300 
                   border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] 
                   hover:from-zinc-200 hover:to-zinc-400 transition select-none flex items-center justify-center"
      >
        Go back Home
      </Link>
    </div>
  );
};

export default NotFound;
