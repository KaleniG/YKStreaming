type Props = {
  title: string;
  children: React.ReactNode;
  onSubmit: (e: React.FormEvent) => void;
};

export const AuthForm: React.FC<Props> = ({ title, children, onSubmit }) => {
  return (
    <div className="flex flex-col h-full items-center flex-1 bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100 pt-16">
      <form
        onSubmit={onSubmit}
        className="w-80 bg-gradient-to-b from-zinc-100 to-zinc-200 rounded-lg shadow-[0_2px_6px_rgba(0,0,0,0.15)] p-6 border border-zinc-400"
      >
        <h2 className="text-2xl font-semibold mb-6 text-center text-zinc-700 select-none">
          {title}
        </h2>
        {children}
      </form>
    </div>
  );
};
