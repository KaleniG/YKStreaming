import * as React from "react";
import { Navigate } from "react-router-dom";
import useAuth from "./AuthContext";

const ProtectedRoute: React.FC<React.PropsWithChildren> = ({ children }) => {
  const auth = useAuth();

  if (!auth.isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
};

export default ProtectedRoute;
