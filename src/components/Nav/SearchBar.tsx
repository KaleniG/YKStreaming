import * as React from "react";

const SearchBar: React.FC = () => {
  return (
    <div className="relative w-full max-w-md">
      <input
        type="text"
        placeholder="Search streams..."
        className="caret-zinc-500 selection:bg-zinc-300 selection:text-black w-full h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border border-zinc-400 px-3 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1 focus:ring-zinc-500"
        autoCorrect="off"
        autoCapitalize="off"
        spellCheck={false}
      />
    </div>
  );
};

export default React.memo(SearchBar);
