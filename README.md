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
cp .env.example .env
go mod download
go run cmd/main.go