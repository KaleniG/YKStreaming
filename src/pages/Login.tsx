import * as React from "react";
import axios from "axios";
import { Link, useNavigate } from "react-router-dom";
import useAuth from "../core/AuthContext";

import { AuthForm } from "../components/AuthForm";
import { FormInput } from "../components/FormInput";
import { PasswordInput } from "../components/PasswordInput";

const Login: React.FC = () => {
  const auth = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [rememberMe, setRememberMe] = React.useState(false);
  const [loading, setLoading] = React.useState(false);
  const [invalidField, setInvalidField] = React.useState<
    "email" | "password" | null
  >(null);

  const emailRef = React.useRef<HTMLInputElement>(null);
  const passwordRef = React.useRef<HTMLInputElement>(null);

  React.useEffect(() => {
    if (auth.isAuthenticated) {
      navigate("/");
    }
  }, [auth, navigate]);

  const validate = () => {
    if (!email) {
      setInvalidField("email");
      emailRef.current?.focus();
      return false;
    }
    if (!password) {
      setInvalidField("password");
      passwordRef.current?.focus();
      return false;
    }
    setInvalidField(null);
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    setLoading(true);

    try {
      const res = await axios.post(
        "http://localhost/api/auth/login",
        { email, password, remember_me: rememberMe },
        { withCredentials: true }
      );

      if (res.status === 200) {
        auth.setAuthenticated(true);
        navigate("/");
      }
    } catch (err: any) {
      const param = err?.response?.data?.param;
      setInvalidField(param ?? null);

      if (param === "email") emailRef.current?.focus();
      if (param === "password") passwordRef.current?.focus();

      auth.setAuthenticated(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthForm title="Login" onSubmit={handleSubmit}>
      <FormInput
        id="email"
        label="Email"
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        invalid={invalidField === "email"}
        inputRef={emailRef}
      />

      <PasswordInput
        id="password"
        label="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        invalid={invalidField === "password"}
        inputRef={passwordRef}
      />

      {/* Remember Me */}
      <label className="flex items-center mb-4 cursor-pointer select-none mt-4">
        <input
          type="checkbox"
          checked={rememberMe}
          onChange={(e) => setRememberMe(e.target.checked)}
          className="mr-2 mt-1 accent-zinc-700"
        />
        <span className="text-zinc-700 text-sm">Remember me</span>
      </label>

      <button
        type="submit"
        disabled={loading}
        className="w-full h-8 rounded bg-gradient-to-b from-zinc-100 to-zinc-300
                   border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)]
                   hover:from-zinc-200 hover:to-zinc-400 
                   font-semibold transition flex justify-center items-center
                   select-none text-zinc-700"
      >
        {loading ? "Logging in..." : "Login"}
      </button>

      <p className="text-sm text-zinc-600 mt-4 text-center select-none">
        Don&apos;t have an account?{" "}
        <Link
          to="/signup"
          className="text-zinc-700 hover:underline font-semibold"
        >
          Sign up
        </Link>
      </p>
    </AuthForm>
  );
};

export default Login;
