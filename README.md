# DevBrain

DevBrain turns the technical articles you read into a personal knowledge graph. Import an article, extract its concepts and relationships, then explore what you know and identify what to learn next.

## Features

- Import and save blog posts and articles from a URL.
- Use Gemini to extract concepts and relationships from saved resources.
- Explore the resulting knowledge graph and inspect individual concepts.
- Compare learned concepts with learning roadmaps for Go, Docker, Kubernetes, system design, databases, and JavaScript.
- Search your saved resources with natural-language questions.
- Ask the AI mentor questions using context from your knowledge graph.
- Review resource and learning analytics on the dashboard.
- Register and sign in with JWT-based authentication.

## How It Works

1. Create an account and sign in.
2. Import an article from the Resources page.
3. Select **Extract concepts** on the saved article. DevBrain extracts concepts and then attempts to extract their relationships.
4. Explore the graph, search your resources, ask the mentor, or review learning gaps.

## Technology

- **Frontend:** Next.js 16, React 19, TypeScript, Tailwind CSS, React Flow
- **Backend:** Go, Gin, GORM
- **Database:** PostgreSQL
- **AI:** Google Gemini API

## Requirements

- Go 1.26.3 or later
- Node.js 20.9 or later and npm
- A PostgreSQL database
- A Google Gemini API key

The backend runs GORM migrations at startup. The database must be reachable before starting the server.

## Run Locally

### 1. Configure the backend

In PowerShell, from the repository root:

```powershell
cd backend
Copy-Item .env.example .env
```

Edit `backend/.env` and set the following values:

| Variable | Required | Description |
| --- | --- | --- |
| `DATABASE_URL` | Yes | PostgreSQL connection string |
| `JWT_SECRET` | Yes | Long, random secret used to sign authentication tokens |
| `GEMINI_API_KEY` | Yes | API key from Google AI Studio |
| `FRONTEND_URL` | No | Additional allowed browser origin; defaults to `http://localhost:3000` |

Keep real credentials in your local `.env` file; do not commit them.

### 2. Start the backend

In a terminal from the `backend` directory:

```powershell
go mod download
go run ./cmd
```

The API listens at `http://localhost:8080`.

### 3. Start the frontend

Open another terminal from the repository root:

```powershell
cd frontend
npm ci
```

The frontend defaults to `http://localhost:8080/api` for its backend API, so no frontend environment file is needed for local development. To override it, copy `frontend/.env.local.example` to `frontend/.env.local` and set `NEXT_PUBLIC_API_URL` (include the `/api` suffix).

Start Next.js:

```powershell
npm run dev
```

Open [http://localhost:3000](http://localhost:3000).

## Useful Commands

Run backend checks from `backend/`:

```powershell
go test ./...
```

Run frontend checks from `frontend/`:

```powershell
npm run lint
npm run build
```

## API Overview

Registration, login, and concept detail are public. The other routes require a bearer token:

| Method | Route | Purpose |
| --- | --- | --- |
| `POST` | `/api/auth/register` | Create an account |
| `POST` | `/api/auth/login` | Sign in |
| `GET` | `/api/concepts/:name` | Get concept details |
| `GET` | `/api/me` | Get the authenticated user |
| `GET`, `POST` | `/api/resources` | List and create resources |
| `GET` | `/api/resources/:id` | Get a resource |
| `POST` | `/api/resources/blog` | Import an article by URL |
| `POST` | `/api/concepts/extract` | Extract concepts from a resource |
| `POST` | `/api/relationships/extract` | Extract resource relationships |
| `GET` | `/api/graph` | Get the user's knowledge graph |
| `POST` | `/api/search` | Search saved resources |
| `POST` | `/api/chat` | Ask the AI mentor |
| `GET` | `/api/gaps`, `/api/gaps/domains` | Get learning gaps and domains |
| `GET` | `/api/analytics` | Get learning analytics |

## Troubleshooting

- **Backend exits while connecting:** verify `DATABASE_URL`, database availability, and network access from the machine running Go.
- **Gemini requests fail:** verify `GEMINI_API_KEY` and that the key is enabled for the Gemini API.
- **Browser cannot reach the API:** confirm the backend is running on port `8080`. If the frontend is hosted elsewhere, set `NEXT_PUBLIC_API_URL` and configure the backend's `FRONTEND_URL` to allow that origin.
- **Mentor has little context:** import a resource and run concept extraction first; mentor answers are grounded in the available graph data.

## Project Structure

```text
backend/
  cmd/                 Server entry point
  internal/            API handlers, services, models, and database setup
frontend/
  src/app/             Next.js pages
  src/components/      Shared UI components
  src/lib/             API client and frontend utilities
```