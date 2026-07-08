"use client";

import { useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import { useRouter } from "next/navigation";
import { Button, Input } from "@/components/ui";

export default function RegisterPage() {
  const router = useRouter();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  async function handleRegister(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await api.post("/auth/register", {
        name,
        email,
        password,
      });

      localStorage.setItem("token", response.data.token);

      router.push("/dashboard");
    } catch (err: any) {
      setError(
        err.response?.data?.error ?? "Couldn't create your account."
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
          <h1 className="text-lg font-semibold mb-1">Create your account</h1>
          <p className="text-sm text-ink-soft mb-6">
            Free, forever. No credit card involved.
          </p>

          <form onSubmit={handleRegister} className="flex flex-col gap-3">
            <label className="flex flex-col gap-1.5">
              <span className="font-label text-xs text-ink-soft">name</span>
              <Input
                required
                placeholder="Ada Lovelace"
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </label>

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
                minLength={6}
                placeholder="At least 6 characters"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
            </label>

            {error && <p className="text-sm text-danger">{error}</p>}

            <Button type="submit" disabled={loading} className="w-full mt-2">
              {loading ? "Creating account…" : "Create account"}
            </Button>
          </form>
        </div>

        <p className="text-center text-sm text-ink-soft mt-5">
          Already have an account?{" "}
          <Link href="/login" className="text-accent font-medium">
            Log in
          </Link>
        </p>
      </div>
    </div>
  );
}
