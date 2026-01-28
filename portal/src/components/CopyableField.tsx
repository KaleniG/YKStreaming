import * as React from "react";
import { IoCopy, IoCopyOutline } from "react-icons/io5";

type Props = {
  value: string;
  className: string;
};

export const CopyableField: React.FC<Props> = ({ value, className }) => {
  const [copyed, setCopyed] = React.useState<boolean>(false);

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopyed(true);
    setTimeout(() => setCopyed(false), 1500);
  };

  const fieldBaseStyle = "caret-zinc-500 selection:bg-zinc-300 selection:text-black h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border px-3 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1 border-zinc-400 rounded-l-md";
  return (
    <div className="relative">
      <input
        type="text"
        value={value}
        disabled
        className={`${fieldBaseStyle} ${className}`}
      />
      <button
        onClick={() => copyToClipboard(value)}
        className={`absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700 ${copyed ? "focus:cursor-default" : ""}`}
      >
        {copyed ? (
          <IoCopy size={18} />
        ) : (
          <IoCopyOutline size={18} />
        )}
      </button>
    </div >
  );
};
