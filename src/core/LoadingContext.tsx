import * as React from "react";

interface PageLoadingContext {
  isLoading: boolean;
  setLoading: (value: boolean) => void;
}

const pageLoadingContext = React.createContext<PageLoadingContext | null>(null);

export const PageLoadingProvider: React.FC<React.PropsWithChildren> = ({
  children,
}) => {
  const [isLoading, setIsLoading] = React.useState(false);

  const ctx = React.useMemo(() => {
    return { isLoading, setLoading: setIsLoading };
  }, [isLoading]);

  return (
    <pageLoadingContext.Provider value={ctx}>
      {children}
    </pageLoadingContext.Provider>
  );
};

const usePageLoading = (): PageLoadingContext => {
  const context = React.useContext(pageLoadingContext);
  if (!context)
    throw new Error("useLoading must be used inside PageLoadingProvider");
  return context;
};

export default usePageLoading;
