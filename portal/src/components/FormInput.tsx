import * as React from "react";

type Props = {
  id: string;
  label: string;
  type?: string;
  value: string;
  onChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  placeholder?: string;
  invalid?: boolean;
  inputRef?: React.RefObject<HTMLInputElement>;
  extraRightPadding?: boolean;
};

const baseStyle =
  "caret-zinc-500 selection:bg-zinc-300 selection:text-black w-full h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border px-3 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1";

export const FormInput: React.FC<Props> = ({
  id,
  label,
  type = "text",
  value,
  onChange,
  placeholder,
  invalid,
  inputRef,
  extraRightPadding,
}) => {
  return (
    <>
      <label
        htmlFor={id}
        className="block text-zinc-700 mb-2 font-medium select-none mt-4"
      >
        {label}
      </label>
      <input
        id={id}
        ref={inputRef}
        type={type}
        value={value}
        onChange={onChange}
        placeholder={placeholder}
        className={`${baseStyle} ${invalid
          ? "border-red-600 focus:ring-red-500"
          : "border-zinc-400 focus:ring-zinc-500"
          } ${extraRightPadding ? "pr-10" : ""}`}
      />
    </>
  );
};
