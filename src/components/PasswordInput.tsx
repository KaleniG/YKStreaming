import * as React from "react";
import { AiOutlineEye, AiOutlineEyeInvisible } from "react-icons/ai";

type Props = {
  id: string;
  label: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  invalid?: boolean;
  inputRef?: React.RefObject<HTMLInputElement>;
};

const baseStyle =
  "caret-zinc-500 selection:bg-zinc-300 selection:text-black w-full h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border px-3 pr-10 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1";

export const PasswordInput: React.FC<Props> = ({
  id,
  label,
  value,
  onChange,
  invalid,
  inputRef,
}) => {
  const [show, setShow] = React.useState(false);

  return (
    <>
      <label
        htmlFor={id}
        className="block text-zinc-700 mb-2 font-medium select-none mt-4"
      >
        {label}
      </label>
      <div className="relative">
        <input
          id={id}
          ref={inputRef}
          type={show ? "text" : "password"}
          value={value}
          onChange={onChange}
          className={`${baseStyle} ${invalid
            ? "border-red-600 focus:ring-red-500"
            : "border-zinc-400 focus:ring-zinc-500"
            }`}
        />
        <button
          type="button"
          onClick={() => setShow((v) => !v)}
          className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
        >
          {show ? <AiOutlineEyeInvisible size={20} /> : <AiOutlineEye size={20} />}
        </button>
      </div>
    </>
  );
};
