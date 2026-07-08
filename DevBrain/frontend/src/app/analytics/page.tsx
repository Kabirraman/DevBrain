"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import AppShell from "@/components/AppShell";
import { Card, PageHeader, ProgressBar, EmptyState } from "@/components/ui";

interface Analytics {
  total_resources: number;
  total_concepts: number;
  most_learned_topics: string[];
  knowledge_coverage: { domain: string; completion: number }[];
  missing_concepts: string[];
}

export default function AnalyticsPage() {
  const [data, setData] = useState<Analytics | null>(null);

  useEffect(() => {
    api.get("/analytics").then((r) => setData(r.data)).catch(console.error);
  }, []);

  if (!data) {
    return (
      <AppShell>
        <p className="text-sm text-ink-soft">Loading…</p>
      </AppShell>
    );
  }

  return (
    <AppShell>
      <PageHeader
        eyebrow="analytics"
        title="Your learning analytics"
        description="A birds-eye view of everything you've imported and learned so far."
      />

      <div className="grid sm:grid-cols-2 gap-4 mb-8">
        <Card>
          <div className="font-label text-xs text-ink-soft mb-2">
            resources imported
          </div>
          <div className="text-3xl font-semibold">{data.total_resources}</div>
        </Card>
        <Card>
          <div className="font-label text-xs text-ink-soft mb-2">
            concepts learned
          </div>
          <div className="text-3xl font-semibold">{data.total_concepts}</div>
        </Card>
      </div>

      <Card className="mb-6">
        <div className="font-label text-xs text-ink-soft mb-4">
          knowledge coverage by domain
        </div>
        {data.knowledge_coverage.length === 0 ? (
          <p className="text-sm text-ink-soft">
            Extract concepts from a resource to see coverage here.
          </p>
        ) : (
          <div className="flex flex-col gap-4">
            {data.knowledge_coverage.map((c) => (
              <div key={c.domain}>
                <div className="flex justify-between text-sm mb-1.5">
                  <span>{c.domain}</span>
                  <span className="text-ink-soft">{c.completion}%</span>
                </div>
                <ProgressBar value={c.completion} />
              </div>
            ))}
          </div>
        )}
      </Card>

      <div className="grid sm:grid-cols-2 gap-4">
        <Card>
          <div className="font-label text-xs text-ink-soft mb-3">
            most learned topics
          </div>
          {data.most_learned_topics.length === 0 ? (
            <p className="text-sm text-ink-soft">Nothing yet.</p>
          ) : (
            <ul className="flex flex-col gap-1.5">
              {data.most_learned_topics.map((t) => (
                <li key={t} className="text-sm capitalize">
                  {t}
                </li>
              ))}
            </ul>
          )}
        </Card>

        <Card>
          <div className="font-label text-xs text-ink-soft mb-3">
            commonly missing next
          </div>
          {data.missing_concepts.length === 0 ? (
            <p className="text-sm text-ink-soft">You&apos;re all caught up.</p>
          ) : (
            <ul className="flex flex-col gap-1.5">
              {data.missing_concepts.map((t) => (
                <li key={t} className="text-sm capitalize text-ink-soft">
                  {t}
                </li>
              ))}
            </ul>
          )}
        </Card>
      </div>
    </AppShell>
  );
}
