# DevBrain

An AI-powered knowledge graph for developers. Import blog posts and articles,
DevBrain extracts the concepts and relationships automatically, and builds
a personal graph of what you know — plus what you're missing.

## Features

- **Auth** — JWT-based register/login
- **Resource ingestion** — import any blog post by URL; content is scraped
  and stored
- **Concept & relationship extraction** — Gemini reads each resource and
  extracts technical concepts + how they relate, feeding a knowledge graph
- **Knowledge graph explorer** — interactive graph view (React Flow) of
  every concept and relationship you've built
- **Learning gap detection** — compares what you've learned against curated
  domain roadmaps (Go, Docker, Kubernetes, System Design, Databases,
  JavaScript) and shows exactly what's missing
- **Resource search** — ask a question in plain English, get an answer
  synthesized from your own imported resources
- **AI mentor chat** — ask questions answered strictly from your knowledge
  graph
- **Analytics dashboard** — resource counts, most-learned topics, coverage
  per domain, and what to learn next

## Stack

- **Backend**: Go, Gin, GORM, PostgreSQL, Gemini API (`gemini-2.5-flash`)
- **Frontend**: Next.js 16, React 19, Tailwind CSS, React Flow

## Local setup

### Backend

```bash
cd backend
cp .env.example .env   # fill in DATABASE_URL, JWT_SECRET, GEMINI_API_KEY
go mod download
go run cmd/main.go
```

Runs on `http://localhost:8080`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Runs on `http://localhost:3000`. It talks to `http://localhost:8080/api` by
default — no `.env.local` needed for local dev.

## Deploying for $0

This whole stack fits comfortably inside free tiers.

1. **Database — [Supabase](https://supabase.com) (free)**
   Create a project, then copy the connection string (Project Settings →
   Database → Connection string → **Session pooler**, so it works from a
   serverless/free host). Use it as `DATABASE_URL`.

2. **Backend — [Render](https://render.com) (free web service)**
   New → Web Service → connect this repo → root directory `backend` →
   Render will detect the `Dockerfile` automatically. Add environment
   variables: `DATABASE_URL`, `JWT_SECRET`, `GEMINI_API_KEY`, and
   `FRONTEND_URL` (your Vercel URL, set after step 3). Note: Render's free
   tier spins down after inactivity, so the first request after a while
   takes ~30-50s to wake up.

3. **Frontend — [Vercel](https://vercel.com) (free)**
   Import this repo → root directory `frontend` → add environment variable
   `NEXT_PUBLIC_API_URL` = `https://your-render-app.onrender.com/api` →
   Deploy.

4. **AI — [Google AI Studio](https://aistudio.google.com/apikey) (free)**
   Generate a Gemini API key on the free tier (generous daily quota for
   `gemini-2.5-flash`) and use it as `GEMINI_API_KEY` in step 2.

Once deployed, go back to Render and set `FRONTEND_URL` to your live Vercel
URL so CORS allows it, then redeploy the backend.

## Notes on design choices

- Concept extraction and search deliberately use lightweight keyword/LLM
  techniques rather than a vector database — this keeps the entire stack on
  free tiers with zero paid infrastructure (no managed vector DB, no
  embeddings API billing). If you want to extend this into full vector
  search later, Supabase supports the `pgvector` extension out of the box.
- Domain roadmaps for gap detection are defined in
  `backend/internal/gaps/roadmaps.go` as plain Go maps — easy to extend
  with new domains without touching the database.
