import * as React from "react";
import { Routes, Route } from "react-router-dom";

import Home from "../pages/Home";
import Login from "../pages/Login";
import Signup from "../pages/Signup";
import UserStreams from "../pages/UserStreams";
import Stream from "../pages/Stream";
import NotFound from "../pages/NotFound";

const AuthRouter: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<Signup />} />
      <Route path="/user-stream" element={<UserStreams />} />
      <Route path="/stream/:streamKey" element={<Stream />} />
      <Route path="*" element={<NotFound />} />
    </Routes>
  );
};

export default AuthRouter;
