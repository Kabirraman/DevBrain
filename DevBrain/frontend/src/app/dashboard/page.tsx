"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import AppShell from "@/components/AppShell";
import { Card, PageHeader, ProgressBar } from "@/components/ui";
import { Library, Network, Target, ArrowRight } from "lucide-react";

interface Me {
  id: string;
  name: string;
  email: string;
}

interface Analytics {
  total_resources: number;
  total_concepts: number;
  most_learned_topics: string[];
  knowledge_coverage: { domain: string; completion: number }[];
}

export default function DashboardPage() {
  const [user, setUser] = useState<Me | null>(null);
  const [analytics, setAnalytics] = useState<Analytics | null>(null);

  useEffect(() => {
    api.get("/me").then((r) => setUser(r.data)).catch(console.error);
    api
      .get("/analytics")
      .then((r) => setAnalytics(r.data))
      .catch(console.error);
  }, []);

  return (
    <AppShell>
      <PageHeader
        eyebrow="dashboard"
        title={user ? `Welcome back, ${user.name.split(" ")[0]}` : "Welcome"}
        description="Here's a snapshot of what you've built so far."
      />

      <div className="grid sm:grid-cols-3 gap-4 mb-8">
        <Card>
          <div className="font-label text-xs text-ink-soft mb-2">
            resources
          </div>
          <div className="text-3xl font-semibold">
            {analytics?.total_resources ?? "–"}
          </div>
        </Card>
        <Card>
          <div className="font-label text-xs text-ink-soft mb-2">
            concepts known
          </div>
          <div className="text-3xl font-semibold">
            {analytics?.total_concepts ?? "–"}
          </div>
        </Card>
        <Card>
          <div className="font-label text-xs text-ink-soft mb-2">
            domains tracked
          </div>
          <div className="text-3xl font-semibold">
            {analytics?.knowledge_coverage.length ?? "–"}
          </div>
        </Card>
      </div>

      {analytics && analytics.knowledge_coverage.length > 0 && (
        <Card className="mb-8">
          <div className="font-label text-xs text-ink-soft mb-4">
            knowledge coverage
          </div>
          <div className="flex flex-col gap-4">
            {analytics.knowledge_coverage.map((c) => (
              <div key={c.domain}>
                <div className="flex justify-between text-sm mb-1.5">
                  <span>{c.domain}</span>
                  <span className="text-ink-soft">{c.completion}%</span>
                </div>
                <ProgressBar value={c.completion} />
              </div>
            ))}
          </div>
        </Card>
      )}

      <div className="grid sm:grid-cols-3 gap-4">
        <QuickLink
          href="/resources"
          icon={Library}
          label="Add a resource"
          text="Import a blog post to start building your graph."
        />
        <QuickLink
          href="/graph"
          icon={Network}
          label="Explore your graph"
          text="See how concepts connect to each other."
        />
        <QuickLink
          href="/gaps"
          icon={Target}
          label="Find your gaps"
          text="Compare what you know against a roadmap."
        />
      </div>
    </AppShell>
  );
}

function QuickLink({
  href,
  icon: Icon,
  label,
  text,
}: {
  href: string;
  icon: any;
  label: string;
  text: string;
}) {
  return (
    <Link
      href={href}
      className="group border border-line rounded-lg p-5 hover:border-ink transition-colors"
    >
      <Icon size={18} className="text-accent mb-3" strokeWidth={2} />
      <div className="flex items-center gap-1.5 text-sm font-medium">
        {label}
        <ArrowRight
          size={14}
          className="opacity-0 group-hover:opacity-100 transition-opacity"
        />
      </div>
      <p className="text-sm text-ink-soft mt-1">{text}</p>
    </Link>
  );
}
