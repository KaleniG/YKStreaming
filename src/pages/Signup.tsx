import * as React from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import { AiOutlineEye, AiOutlineEyeInvisible } from "react-icons/ai";

import useAuth from "../core/AuthContext";

const Signup: React.FC = () => {
  const auth = useAuth();
  const navigate = useNavigate();

  const [username, setUsername] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [showPassword, setShowPassword] = React.useState(false);
  const [repeatPassword, setRepeatPassword] = React.useState("");
  const [loading, setLoading] = React.useState(false);
  const [invalidField, setInvalidField] = React.useState<
    "email" | "password" | "repeatPassword" | null
  >(null);

  const emailRef = React.useRef<HTMLInputElement>(null);
  const passwordRef = React.useRef<HTMLInputElement>(null);
  const repeatPasswordRef = React.useRef<HTMLInputElement>(null);

  React.useEffect(() => {
    if (auth.isAuthenticated) {
      navigate("/");
    }
  }, [auth, navigate]);

  const validate = (): boolean => {
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
    if (!repeatPassword) {
      setInvalidField("repeatPassword");
      repeatPasswordRef.current?.focus();
      return false;
    }
    if (password !== repeatPassword) {
      setInvalidField("repeatPassword");
      repeatPasswordRef.current?.focus();
      return false;
    }
    setInvalidField(null);
    return true;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validate()) return;

    setInvalidField(null);
    setLoading(true);

    try {
      const res = await axios.post(
        "http://localhost/api/auth/signup",
        { name: username, email, password },
        { withCredentials: true }
      );

      if (res.status == 200) {
        auth.setAuthenticated(true);
        navigate("/");
      }
    } catch (err: any) {
      if (err?.response?.data) {
        console.warn(err.response.data.error);
        auth.setAuthenticated(false);
        setInvalidField(err.response.data.param);
        if (err.response.data.param == "email") {
          emailRef.current?.focus();
        }
      }
    } finally {
      setLoading(false);
    }
  };

  const inputFieldBaseStyle =
    "caret-zinc-500 selection:bg-zinc-300 selection:text-black w-full h-9 rounded-md bg-gradient-to-b from-white to-zinc-200 border px-3 text-sm shadow-[inset_0_1px_2px_rgba(0,0,0,0.2)] focus:outline-none focus:ring-1";

  return (
    <div className="flex flex-col h-full items-center flex-1 bg-gradient-to-b from-zinc-100 via-zinc-50 to-zinc-100 pt-16">
      <form
        onSubmit={handleSubmit}
        className="w-80 bg-gradient-to-b from-zinc-100 to-zinc-200 rounded-lg shadow-[0_2px_6px_rgba(0,0,0,0.15)] p-6 border border-zinc-400"
      >
        <h2 className="text-2xl font-semibold mb-6 text-center text-zinc-700 select-none">
          Sign up
        </h2>

        {/* Username */}
        <label
          className="block text-zinc-700 mb-2 font-medium select-none mt-4"
          htmlFor="username"
        >
          Username
        </label>
        <input
          id="username"
          type="text"
          placeholder="Enter your username"
          className={`${inputFieldBaseStyle} border-zinc-400 focus:ring-zinc-500`}
          value={username}
          onChange={(e) => setUsername(e.target.value)}
        />

        {/* Email */}
        <label
          className="block text-zinc-700 mb-2 font-medium select-none mt-4"
          htmlFor="email"
        >
          Email
        </label>
        <input
          id="email"
          type="email"
          ref={emailRef}
          placeholder="Enter your email"
          className={`${inputFieldBaseStyle} ${invalidField === "email"
            ? "border-red-600 focus:ring-red-500"
            : "border-zinc-400 focus:ring-zinc-500"
            }`}
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          autoCorrect="off"
          autoCapitalize="off"
          spellCheck={false}
        />

        {/* Password */}
        <label
          className="block text-zinc-700 mb-2 font-medium select-none mt-4"
          htmlFor="password"
        >
          Password
        </label>
        <div className="relative">
          <input
            id="password"
            type={showPassword ? "text" : "password"}
            ref={passwordRef}
            placeholder="Enter your password"
            className={`${inputFieldBaseStyle} ${invalidField === "password"
              ? "border-red-600 focus:ring-red-500"
              : "border-zinc-400 focus:ring-zinc-500"
              }`}
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <button
            type="button"
            onClick={() => setShowPassword((prev) => !prev)}
            className="absolute right-2 top-1/2 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
          >
            {showPassword ? (
              <AiOutlineEyeInvisible size={20} />
            ) : (
              <AiOutlineEye size={20} />
            )}
          </button>
        </div>

        {/* Repeat Password */}
        <label
          className="block text-zinc-700 mb-2 font-medium select-none mt-4"
          htmlFor="repeatPassword"
        >
          Repeat Password
        </label>
        <div className="relative">
          <input
            id="repeatPassword"
            type={showPassword ? "text" : "password"}
            ref={repeatPasswordRef}
            placeholder="Repeat your password"
            className={`${inputFieldBaseStyle} ${invalidField === "repeatPassword"
              ? "border-red-600 focus:ring-red-500"
              : "border-zinc-400 focus:ring-zinc-500"
              } mb-4`}
            value={repeatPassword}
            onChange={(e) => setRepeatPassword(e.target.value)}
          />
          <button
            type="button"
            onClick={() => setShowPassword((prev) => !prev)}
            className="absolute right-2 top-1/3 -translate-y-1/2 text-zinc-500 hover:text-zinc-700"
          >
            {showPassword ? (
              <AiOutlineEyeInvisible size={20} />
            ) : (
              <AiOutlineEye size={20} />
            )}
          </button>
        </div>

        <button
          type="submit"
          disabled={loading}
          className="w-full h-8 rounded bg-gradient-to-b from-zinc-100 to-zinc-300
                     border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)]
                     hover:from-zinc-200 hover:to-zinc-400 
                     font-semibold transition flex justify-center items-center select-none text-zinc-700 mt-4"
        >
          {loading ? "Signing up..." : "Sign up"}
        </button>
      </form>
    </div>
  );
};

export default Signup;
