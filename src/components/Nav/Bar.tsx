import * as React from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

import useAuth from "../../core/AuthContext";
import usePageLoading from "../../core/LoadingContext";
import Logo from "./Logo";
import SearchBar from "./SearchBar";
import Button from "./Button";

const Bar: React.FC = () => {
  const auth = useAuth();
  const page = usePageLoading();
  const navigate = useNavigate();

  const handleLogout = async () => {
    page.setLoading(true);
    try {
      const res = await axios.post(
        "http://localhost/api/user/logout",
        {},
        { withCredentials: true }
      );

      if (res.status == 200) {
        auth.setAuthenticated(false);
        navigate("/");
      }
    } catch (err: any) {
      if (err?.response?.data) {
        console.warn(err?.response?.data.error)
      }
      if (err.response?.status == 401) {
        auth.setAuthenticated(false)
        navigate("/")
      }
    } finally {
      page.setLoading(false);
    }
  };

  return (
    <nav className="w-full h-14 border-b border-zinc-400 bg-gradient-to-b from-zinc-200 via-zinc-100 to-zinc-300 shadow-[inset_0_1px_0_rgba(255,255,255,0.9)]">
      <div className="w-full h-full px-4 flex items-center">
        <div className="flex items-center gap-6">
          <Logo />
        </div>

        <div className="flex-1 flex justify-center">{/*<SearchBar />*/}</div>

        <div className="flex items-center gap-3">
          {auth.isAuthenticated ? (
            <>
              <Button redirect="/user-stream">Stream</Button>
              <div className="h-6 w-px bg-zinc-400/70 self-center"></div>
              <Button onClick={handleLogout}>Logout</Button>
            </>
          ) : (
            <>
              <Button redirect="/login">Login</Button>
              <Button redirect="/signup">Sign up</Button>
            </>
          )}
        </div>
      </div>
    </nav>
  );
};

export default Bar;
