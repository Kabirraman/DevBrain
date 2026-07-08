"use client";

import { useState, useRef, useEffect } from "react";
import { api } from "@/lib/api";
import AppShell from "@/components/AppShell";
import { PageHeader, Input, Button } from "@/components/ui";
import { Send, Bot, User } from "lucide-react";

interface Message {
  role: "user" | "mentor";
  text: string;
  concepts?: string[];
}

export default function MentorPage() {
  const [messages, setMessages] = useState<Message[]>([]);
  const [question, setQuestion] = useState("");
  const [loading, setLoading] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  async function handleSend(e: React.FormEvent) {
    e.preventDefault();
    if (!question.trim()) return;

    const q = question;
    setMessages((prev) => [...prev, { role: "user", text: q }]);
    setQuestion("");
    setLoading(true);

    try {
      const res = await api.post("/chat", { question: q });
      setMessages((prev) => [
        ...prev,
        {
          role: "mentor",
          text: res.data.answer,
          concepts: res.data.concepts_used,
        },
      ]);
    } catch (err) {
      setMessages((prev) => [
        ...prev,
        { role: "mentor", text: "Something went wrong answering that. Try again." },
      ]);
    } finally {
      setLoading(false);
    }
  }

  return (
    <AppShell>
      <PageHeader
        eyebrow="mentor"
        title="Ask your AI mentor"
        description="Answers are grounded only in your knowledge graph — the concepts and relationships you've already extracted."
      />

      <div className="border border-line rounded-lg flex flex-col h-[55vh]">
        <div className="flex-1 overflow-y-auto p-5 flex flex-col gap-4">
          {messages.length === 0 && (
            <p className="text-sm text-ink-soft m-auto text-center max-w-xs">
              Ask something like &ldquo;How does Go handle concurrency?&rdquo;
              — the mentor answers from your graph only.
            </p>
          )}

          {messages.map((m, i) => (
            <div
              key={i}
              className={`flex gap-2.5 max-w-[85%] ${
                m.role === "user" ? "self-end flex-row-reverse" : "self-start"
              }`}
            >
              <div className="w-6 h-6 rounded-full bg-line-soft flex items-center justify-center shrink-0 mt-0.5">
                {m.role === "user" ? (
                  <User size={13} />
                ) : (
                  <Bot size={13} className="text-accent" />
                )}
              </div>
              <div>
                <div
                  className={`text-sm rounded-lg px-3.5 py-2.5 leading-relaxed whitespace-pre-wrap ${
                    m.role === "user"
                      ? "bg-ink text-paper"
                      : "bg-line-soft text-ink"
                  }`}
                >
                  {m.text}
                </div>
                {m.concepts && m.concepts.length > 0 && (
                  <div className="flex flex-wrap gap-1.5 mt-1.5">
                    {m.concepts.map((c) => (
                      <span
                        key={c}
                        className="font-label text-[11px] text-ink-soft border border-line rounded px-1.5 py-0.5"
                      >
                        {c}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
          ))}
          <div ref={bottomRef} />
        </div>

        <form
          onSubmit={handleSend}
          className="border-t border-line p-3 flex gap-2"
        >
          <Input
            placeholder="Ask a question…"
            value={question}
            onChange={(e) => setQuestion(e.target.value)}
            className="flex-1"
          />
          <Button type="submit" disabled={loading}>
            <Send size={14} />
          </Button>
        </form>
      </div>
    </AppShell>
  );
}
