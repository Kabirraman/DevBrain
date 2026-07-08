"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api";
import AppShell from "@/components/AppShell";
import { Card, PageHeader, Button, Input, EmptyState } from "@/components/ui";
import { Link as LinkIcon, Sparkles, Check } from "lucide-react";

interface Resource {
  ID: string;
  Title: string;
  Type: string;
  SourceURL: string;
  Content: string;
  CreatedAt: string;
}

export default function ResourcesPage() {
  const [resources, setResources] = useState<Resource[]>([]);
  const [loading, setLoading] = useState(true);
  const [url, setUrl] = useState("");
  const [importing, setImporting] = useState(false);
  const [error, setError] = useState("");
  const [extractingID, setExtractingID] = useState<string | null>(null);
  const [extractedIDs, setExtractedIDs] = useState<Set<string>>(new Set());

  async function loadResources() {
    try {
      const res = await api.get("/resources");
      setResources(res.data ?? []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadResources();
  }, []);

  async function handleImport(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setImporting(true);

    try {
      await api.post("/resources/blog", { url });
      setUrl("");
      await loadResources();
    } catch (err: any) {
      setError(
        err.response?.data?.error ?? "Couldn't import that URL. Try another one."
      );
    } finally {
      setImporting(false);
    }
  }

  async function handleExtract(resourceID: string) {
    setExtractingID(resourceID);

    try {
      await api.post("/concepts/extract", { resource_id: resourceID });

      const relRes = await api
        .post("/relationships/extract", { resource_id: resourceID })
        .catch(() => null);

      setExtractedIDs((prev) => new Set(prev).add(resourceID));
    } catch (err) {
      console.error(err);
    } finally {
      setExtractingID(null);
    }
  }

  return (
    <AppShell>
      <PageHeader
        eyebrow="resources"
        title="Your resources"
        description="Import a blog post or article. DevBrain reads it, then extracts the concepts and relationships into your knowledge graph."
      />

      <Card className="mb-8">
        <form onSubmit={handleImport} className="flex gap-2">
          <Input
            required
            type="url"
            placeholder="https://a-blog-post-you-read.com/article"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            className="flex-1"
          />
          <Button type="submit" disabled={importing}>
            <LinkIcon size={14} />
            {importing ? "Importing…" : "Import"}
          </Button>
        </form>
        {error && <p className="text-sm text-danger mt-2">{error}</p>}
      </Card>

      {loading ? (
        <p className="text-sm text-ink-soft">Loading…</p>
      ) : resources.length === 0 ? (
        <EmptyState
          title="No resources yet"
          description="Paste a blog post URL above to import your first resource and start building your graph."
        />
      ) : (
        <div className="flex flex-col gap-3">
          {resources.map((r) => (
            <Card key={r.ID} className="flex items-center justify-between gap-4">
              <div className="min-w-0">
                <div className="font-medium text-sm truncate">{r.Title}</div>
                <a
                  href={r.SourceURL}
                  target="_blank"
                  rel="noreferrer"
                  className="text-xs text-ink-soft hover:text-accent truncate block mt-0.5"
                >
                  {r.SourceURL}
                </a>
              </div>

              <Button
                variant="secondary"
                onClick={() => handleExtract(r.ID)}
                disabled={extractingID === r.ID}
                className="shrink-0"
              >
                {extractedIDs.has(r.ID) ? (
                  <>
                    <Check size={14} /> Extracted
                  </>
                ) : (
                  <>
                    <Sparkles size={14} />
                    {extractingID === r.ID ? "Extracting…" : "Extract concepts"}
                  </>
                )}
              </Button>
            </Card>
          ))}
        </div>
      )}
    </AppShell>
  );
}
