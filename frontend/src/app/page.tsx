import Link from "next/link";

export default function Home() {
  return (
    <main className="min-h-screen flex flex-col items-center justify-center gap-4">
      <h1 className="text-5xl font-bold">
        DevBrain
      </h1>

      <Link href="/login">
        Login
      </Link>

      <Link href="/register">
        Register
      </Link>

      <Link href="/dashboard">
        Dashboard
      </Link>

      <Link href="/graph">
        Graph
      </Link>
    </main>
  );
}