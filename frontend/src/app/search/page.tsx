"use client";

import { useState } from "react";
import { api } from "@/lib/api";
import AppShell from "@/components/AppShell";
import { Card, PageHeader, Button, Input, EmptyState } from "@/components/ui";
import { Search as SearchIcon } from "lucide-react";

interface SearchResult {
  resource_id: string;
  title: string;
  snippet: string;
  score: number;
}

interface SearchResponse {
  answer: string;
  results: SearchResult[];
}

export default function SearchPage() {
  const [query, setQuery] = useState("");
  const [response, setResponse] = useState<SearchResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const res = await api.post("/search", { query });
      setResponse(res.data);
    } catch (err: any) {
      setError(err.response?.data?.error ?? "Search failed. Try again.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <AppShell>
      <PageHeader
        eyebrow="search"
        title="Search your resources"
        description="Ask a question in plain language. DevBrain looks across everything you've imported and answers from it directly."
      />

      <form onSubmit={handleSearch} className="flex gap-2 mb-8">
        <Input
          required
          placeholder="How does Go handle concurrency?"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="flex-1"
        />
        <Button type="submit" disabled={loading}>
          <SearchIcon size={14} />
          {loading ? "Searching…" : "Search"}
        </Button>
      </form>

      {error && <p className="text-sm text-danger mb-4">{error}</p>}

      {response && (
        <>
          <Card className="mb-6 bg-accent-soft border-accent/30">
            <div className="font-label text-xs text-accent mb-2">answer</div>
            <p className="text-sm leading-relaxed whitespace-pre-wrap">
              {response.answer}
            </p>
          </Card>

          {response.results.length > 0 && (
            <div className="flex flex-col gap-3">
              <div className="font-label text-xs text-ink-soft">sources</div>
              {response.results.map((r) => (
                <Card key={r.resource_id}>
                  <div className="font-medium text-sm mb-1">{r.title}</div>
                  <p className="text-sm text-ink-soft leading-relaxed">
                    …{r.snippet}…
                  </p>
                </Card>
              ))}
            </div>
          )}
        </>
      )}

      {!response && !loading && (
        <EmptyState
          title="No search yet"
          description="Ask something above — DevBrain will search across every resource you've imported."
        />
      )}
    </AppShell>
  );
}
