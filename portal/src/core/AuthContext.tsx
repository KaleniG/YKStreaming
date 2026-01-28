import * as React from "react";
import axios from "axios";
import usePageLoading from "./LoadingContext";

interface AuthContext {
  isAuthenticated: boolean;
  setAuthenticated: (value: boolean) => void;
}

const authContext = React.createContext<AuthContext | null>(null);

interface AuthCheckResponse {
  user: boolean | { name: string; email: string; };
}

export const AuthProvider: React.FC<React.PropsWithChildren> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = React.useState(false);
  const page = usePageLoading();

  React.useEffect(() => {
    const checkAuth = async () => {
      try {
        const res = await axios.post<AuthCheckResponse>(
          "http://localhost/api/auth/check",
          {},
          { withCredentials: true }
        );

        if (res.data?.user) {
          setIsAuthenticated(true);
        } else {
          setIsAuthenticated(false);
        }

      } catch (err: any) {
        if (err?.response?.data) {
          setIsAuthenticated(false);
          console.warn(err?.response?.data.error)
        }
      } finally {
        page.setLoading(false);
      }

    };

    checkAuth();
  }, []);

  const ctx = React.useMemo(() => {
    return {
      isAuthenticated,
      setAuthenticated: setIsAuthenticated,
    };
  }, [isAuthenticated]);

  return <authContext.Provider value={ctx}>{children}</authContext.Provider>;
};

const useAuth = (): AuthContext => {
  const context = React.useContext(authContext);
  if (!context) throw new Error("useAuth must be used inside AuthProvider");
  return context;
};

export default useAuth;
