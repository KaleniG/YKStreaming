import axios from "axios";
import { useNavigate } from "react-router-dom";
import useAuth from "@/core/AuthContext";

export const useAuthGuard = () => {
  const auth = useAuth();
  const navigate = useNavigate();

  return async () => {
    try {
      await axios.post(
        "http://localhost/api/user/logout",
        {},
        { withCredentials: true }
      );
    } catch {
      /* ignore */
    }

    auth.setAuthenticated(false);
    navigate("/login");
  };
};
