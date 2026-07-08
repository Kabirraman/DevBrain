"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import AppShell from "@/components/AppShell";
import { Card, PageHeader, ProgressBar } from "@/components/ui";
import { Check, X } from "lucide-react";

interface GapResult {
  domain: string;
  completion: number;
  known: string[];
  missing: string[];
}

export default function GapsPage() {
  const [domains, setDomains] = useState<string[]>([]);
  const [domain, setDomain] = useState("Go");
  const [result, setResult] = useState<GapResult | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .get("/gaps/domains")
      .then((r) => setDomains(r.data.domains ?? []))
      .catch(console.error);
  }, []);

  useEffect(() => {
    setLoading(true);
    api
      .get(`/gaps?domain=${encodeURIComponent(domain)}`)
      .then((r) => setResult(r.data))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [domain]);

  return (
    <AppShell>
      <PageHeader
        eyebrow="gaps"
        title="Learning gap detection"
        description="Pick a domain roadmap and see exactly which concepts you've already covered — and which ones are still missing."
      />

      <div className="flex flex-wrap gap-2 mb-8">
        {domains.map((d) => (
          <button
            key={d}
            onClick={() => setDomain(d)}
            className={`font-label text-xs px-3 py-1.5 rounded-md border transition-colors ${
              domain === d
                ? "bg-ink text-paper border-ink"
                : "border-line text-ink-soft hover:border-ink hover:text-ink"
            }`}
          >
            {d}
          </button>
        ))}
      </div>

      {loading || !result ? (
        <p className="text-sm text-ink-soft">Loading…</p>
      ) : (
        <>
          <Card className="mb-6">
            <div className="flex justify-between text-sm mb-2">
              <span className="font-medium">{result.domain} coverage</span>
              <span className="text-ink-soft">{result.completion}%</span>
            </div>
            <ProgressBar value={result.completion} />
          </Card>

          <div className="grid sm:grid-cols-2 gap-4">
            <Card>
              <div className="font-label text-xs text-ink-soft mb-3">
                you know
              </div>
              {result.known.length === 0 ? (
                <p className="text-sm text-ink-soft">
                  Nothing here yet — extract concepts from a resource first.
                </p>
              ) : (
                <ul className="flex flex-col gap-2">
                  {result.known.map((k) => (
                    <li
                      key={k}
                      className="flex items-center gap-2 text-sm capitalize"
                    >
                      <Check size={14} className="text-accent shrink-0" />
                      {k}
                    </li>
                  ))}
                </ul>
              )}
            </Card>

            <Card>
              <div className="font-label text-xs text-ink-soft mb-3">
                still missing
              </div>
              {result.missing.length === 0 ? (
                <p className="text-sm text-ink-soft">
                  Nothing missing — you&apos;ve covered this roadmap.
                </p>
              ) : (
                <ul className="flex flex-col gap-2">
                  {result.missing.map((m) => (
                    <li
                      key={m}
                      className="flex items-center gap-2 text-sm capitalize text-ink-soft"
                    >
                      <X size={14} className="text-danger shrink-0" />
                      {m}
                    </li>
                  ))}
                </ul>
              )}
            </Card>
          </div>
        </>
      )}
    </AppShell>
  );
}
