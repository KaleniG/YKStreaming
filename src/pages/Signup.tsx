import * as React from "react";
import axios from "axios";
import { useNavigate } from "react-router-dom";
import useAuth from "../core/AuthContext";

import { AuthForm } from "../components/AuthForm";
import { FormInput } from "../components/FormInput";
import { PasswordInput } from "../components/PasswordInput";

const Signup: React.FC = () => {
  const auth = useAuth();
  const navigate = useNavigate();

  const [username, setUsername] = React.useState("");
  const [email, setEmail] = React.useState("");
  const [password, setPassword] = React.useState("");
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
    if (!repeatPassword || password !== repeatPassword) {
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

    setLoading(true);

    try {
      const res = await axios.post(
        "http://localhost/api/auth/signup",
        { name: username, email, password },
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
      auth.setAuthenticated(false);
    } finally {
      setLoading(false);
    }
  };

  return (
    <AuthForm title="Sign up" onSubmit={handleSubmit}>
      <FormInput
        id="username"
        label="Username"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
      />

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

      <PasswordInput
        id="repeatPassword"
        label="Repeat Password"
        value={repeatPassword}
        onChange={(e) => setRepeatPassword(e.target.value)}
        invalid={invalidField === "repeatPassword"}
        inputRef={repeatPasswordRef}
      />

      <button
        type="submit"
        disabled={loading}
        className="w-full h-8 mt-4 rounded bg-gradient-to-b from-zinc-100 to-zinc-300
                   border border-zinc-400 shadow-[inset_0_1px_0_rgba(255,255,255,0.8)]
                   hover:from-zinc-200 hover:to-zinc-400 
                   font-semibold transition flex justify-center items-center
                   select-none text-zinc-700"
      >
        {loading ? "Signing up..." : "Sign up"}
      </button>
    </AuthForm>
  );
};

export default Signup;
