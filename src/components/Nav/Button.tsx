import * as React from "react";
import { Link } from "react-router-dom";

interface Props {
  children: React.ReactNode; // better than ReactElement for flexibility
  redirect?: string | null;
  onClick?: () => void;
}

const Button: React.FC<Props> = ({ children, redirect, onClick }) => {
  if (redirect) {
    return (
      <Link
        to={redirect}
        className={
          "h-8 pt-1 px-4 text-sm rounded bg-gradient-to-b from-zinc-100 to-zinc-300 border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)] hover:from-zinc-200 hover:to-zinc-400 transition select-none"
        }
      >
        {children}
      </Link>
    );
  }

  return (
    <button
      onClick={onClick}
      className="h-8 px-4 text-sm rounded bg-gradient-to-b from-zinc-100
       to-zinc-300 border border-zinc-400
        shadow-[inset_0_1px_0_rgba(255,255,255,0.8)]
         hover:from-zinc-200 hover:to-zinc-400 transition select-none"
    >
      {children}
    </button>
  );
};

export default React.memo(Button);
