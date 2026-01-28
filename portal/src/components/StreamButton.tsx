import * as React from "react";

type Props = {
  onClick: () => void;
  children: React.ReactNode;
};

export const StreamButton: React.FC<Props> = ({ onClick, children }) => {
  const buttonStyle =
    "mb-2 full-h px-4 text-sm bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none ml-2 bg-gray-300 text-gray-800 hover:bg-gray-400";

  return (
    <button
      onClick={onClick}
      className={buttonStyle}
    >
      {children}
    </button>
  );
};
