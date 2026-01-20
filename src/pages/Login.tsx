import * as React from "react";
import axios from "axios";
import { useNavigate, Link } from "react-router-dom";
import { AiOutlineEye, AiOutlineEyeInvisible } from "react-icons/ai";
import useAuth from "../core/AuthContext";

const Login: React.FC = () => {
  const statusAuth = useAuth();
  const navigate = useNavigate();

  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
  const [showPassword, setShowPassword] = React.useState(false);
  const [rememberMe, setRememberMe] = React.useState(false);
  const [loading, setLoading] = React.useState(false);
  const [invalidField, setInvalidField] = React.useState<
    "email" | "password" | null
  >(null);

  const emailRef = React.useRef<HTMLInputElement>(null);
  const passwordRef = React.useRef<HTMLInputElement>(null);

  React.useEffect(() => {
    if (statusAuth.isAuthenticated) {
      navigate("/");
    }
  }, [statusAuth, navigate]);

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
        "http://localhost/api/auth/login",
        { email, password, remember_me: rememberMe },
        { withCredentials: true }
      );

      if (res.status == 200) {
        statusAuth.setAuthenticated(true);
        navigate("/");
      }
    } catch (err: any) {
      if (err?.response?.data) {
        console.warn(err.response.data.error);
        statusAuth.setAuthenticated(false);
        setInvalidField(err.response.data.param);
        if (err.response.data.param == "email") {
          emailRef.current?.focus();
        } else if (err.response.data.param == "password") {
          passwordRef.current?.focus();
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
          Login
        </h2>

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

        {/* Password with icon */}
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
              } pr-10`} // extra padding for the icon
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

        {/* Submit */}
        <button
          type="submit"
          disabled={loading}
          className="w-full h-8 rounded bg-gradient-to-b from-zinc-100 to-zinc-300
                     border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)]
                     hover:from-zinc-200 hover:to-zinc-400 
                     font-semibold transition flex justify-center items-center select-none text-zinc-700"
        >
          {loading ? "Logging in..." : "Login"}
        </button>

        {/* Signup link */}
        <p className="text-sm text-zinc-600 mt-4 text-center select-none">
          Don&apos;t have an account?{" "}
          <Link
            to="/signup"
            className="text-zinc-700 hover:underline font-semibold"
          >
            Sign up
          </Link>
        </p>
      </form>
    </div>
  );
};

export default Login;
