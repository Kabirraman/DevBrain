"use client";

import { useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";
import { Button, Input } from "@/components/ui";

export default function LoginPage() {
  const router = useRouter();

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleLogin(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await api.post("/auth/login", {
        email,
        password,
      });

      localStorage.setItem("token", response.data.token);

      router.push("/dashboard");
    } catch (err: any) {
      setError(
        err.response?.data?.error ?? "Couldn't log you in. Check your details."
      );
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-paper text-ink px-6">
      <div className="w-full max-w-sm">
        <Link href="/" className="font-label text-sm block mb-10 text-center">
          &gt; devbrain
        </Link>

        <div className="border border-line rounded-lg p-7">
          <h1 className="text-lg font-semibold mb-1">Welcome back</h1>
          <p className="text-sm text-ink-soft mb-6">
            Log in to pick up where you left off.
          </p>

          <form onSubmit={handleLogin} className="flex flex-col gap-3">
            <label className="flex flex-col gap-1.5">
              <span className="font-label text-xs text-ink-soft">email</span>
              <Input
                type="email"
                required
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </label>

            <label className="flex flex-col gap-1.5">
              <span className="font-label text-xs text-ink-soft">
                password
              </span>
              <Input
                type="password"
                required
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </label>

            {error && <p className="text-sm text-danger">{error}</p>}

            <Button type="submit" disabled={loading} className="w-full mt-2">
              {loading ? "Logging in…" : "Log in"}
            </Button>
          </form>
        </div>

        <p className="text-center text-sm text-ink-soft mt-5">
          New here?{" "}
          <Link href="/register" className="text-accent font-medium">
            Create an account
          </Link>
        </p>
      </div>
    </div>
  );
}
