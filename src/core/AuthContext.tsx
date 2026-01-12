import * as React from "react";
import axios from "axios";
import usePageLoading from "./LoadingContext";

interface AuthContext {
  isAuthenticated: boolean;
  setAuthenticated: (value: boolean) => void;
}

const authContext = React.createContext<AuthContext | null>(null);

export const AuthProvider: React.FC<React.PropsWithChildren> = ({
  children,
}) => {
  const [isAuthenticated, setIsAuthenticated] = React.useState(false);
  const page = usePageLoading();

  React.useEffect(() => {
    const checkAuth = async () => {
      try {
        page.setLoading(true);
        const res = await axios.post<{
          success: boolean;
          logged_in: boolean;
        }>(
          "http://localhost/api/auth_check.php",
          {},
          { withCredentials: true }
        );

        setIsAuthenticated(!!res.data?.logged_in);
      } catch {
        setIsAuthenticated(false);
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
