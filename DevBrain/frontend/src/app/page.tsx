import Link from "next/link";
import { Network, Target, Search, MessageSquare } from "lucide-react";

const FEATURES = [
  {
    icon: Network,
    label: "knowledge graph",
    text: "Every concept you learn becomes a node. Every relationship becomes an edge.",
  },
  {
    icon: Target,
    label: "gap detection",
    text: "See exactly what's missing between what you know and a domain roadmap.",
  },
  {
    icon: Search,
    label: "resource search",
    text: "Ask a question, get an answer pulled straight from what you've already read.",
  },
  {
    icon: MessageSquare,
    label: "ai mentor",
    text: "Explain any concept in terms of what you already understand.",
  },
];

export default function Home() {
  return (
    <main className="min-h-screen bg-paper text-ink flex flex-col">
      <header className="flex items-center justify-between px-8 py-6 max-w-5xl mx-auto w-full">
        <span className="font-label text-sm">&gt; devbrain</span>
        <nav className="flex items-center gap-6 text-sm font-label">
          <Link href="/login" className="text-ink-soft hover:text-ink">
            log in
          </Link>
          <Link
            href="/register"
            className="bg-ink text-paper px-3.5 py-1.5 rounded-md hover:bg-accent transition-colors"
          >
            get started
          </Link>
        </nav>
      </header>

      <section className="flex-1 flex flex-col justify-center max-w-5xl mx-auto w-full px-8 py-20">
        <div className="font-label text-xs text-accent mb-4">
          &gt; ai-powered developer knowledge intelligence
        </div>

        <h1 className="text-5xl sm:text-6xl font-semibold tracking-tight max-w-2xl leading-[1.05]">
          Turn what you read into what you know.
        </h1>

        <p className="text-ink-soft mt-6 max-w-lg text-base leading-relaxed">
          DevBrain reads your blog posts, docs, and notes, then builds a
          personal knowledge graph — so you always know what you&apos;ve
          learned, what&apos;s missing, and what to learn next.
        </p>

        <div className="flex items-center gap-4 mt-8">
          <Link
            href="/register"
            className="bg-ink text-paper px-4 py-2.5 rounded-md text-sm font-medium hover:bg-accent transition-colors"
          >
            Create free account
          </Link>
          <Link
            href="/login"
            className="text-sm font-medium text-ink-soft hover:text-ink"
          >
            I already have one →
          </Link>
        </div>

        <div className="grid sm:grid-cols-2 gap-4 mt-20">
          {FEATURES.map((f) => (
            <div
              key={f.label}
              className="border border-line rounded-lg p-5"
            >
              <f.icon size={18} className="text-accent mb-3" strokeWidth={2} />
              <div className="font-label text-xs text-ink-soft mb-1">
                {f.label}
              </div>
              <p className="text-sm text-ink leading-relaxed">{f.text}</p>
            </div>
          ))}
        </div>
      </section>

      <footer className="text-center text-xs text-ink-soft py-6 font-label">
        built with go, next.js &amp; gemini — zero-cost stack
      </footer>
    </main>
  );
}
