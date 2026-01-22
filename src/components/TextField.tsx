import * as React from "react";

type Props = {
  value: string;
  className: string;
};

export const TextField: React.FC<Props> = ({ value, className }) => {
  const fieldBaseStyle = "caret-zinc-500 selection:bg-zinc-300 selection:text-black h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border px-3 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1 border-zinc-400 rounded-l-md";
  return (
    <input
      type="text"
      value={value}
      disabled
      className={`${fieldBaseStyle} ${className}`}
    />
  );
};
