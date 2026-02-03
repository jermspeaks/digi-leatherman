# Digi Leatherman

A full-stack app for string tools (and later other tool types). The **Go backend** exposes an API; the **React + TypeScript** frontend (Vite) calls it.

## Run the app

### Backend (Go API)

From the repo root:

```bash
cd backend
go run .
```

The API listens on **http://localhost:8100**. All string endpoints use body `{"value": "..."}` and return `{"result": "..."}`.

**String API endpoints:**

- **URL:** `POST /api/string/url-encode`, `POST /api/string/url-decode`
- **Base64:** `POST /api/string/base64-encode`, `POST /api/string/base64-decode`
- **Trim:** `POST /api/string/trim`
- **Case:** `POST /api/string/upper-case`, `POST /api/string/lower-case`, `POST /api/string/capital-case`, `POST /api/string/snake-case`, `POST /api/string/kebab-case`, `POST /api/string/camel-case`, `POST /api/string/pascal-case`, `POST /api/string/sentence-case`

### Frontend (Vite + React)

In another terminal:

```bash
cd frontend
npm install   # if you haven't already
npm run dev
```

Open **http://localhost:5273**. The UI has a **Strings** section in the sidebar with collapsible groups (URL, Base64, Trim, Case). Each tool has its own route (e.g. `/tools/string/url-encode`); enter text, run the action, and see the result. The app uses **Tailwind CSS** and a theme-aware layout with a collapsible sidebar. The frontend uses `VITE_API_URL` (default `http://localhost:8100`); copy `.env.example` to `.env` and change it if your API runs elsewhere. If you see a CORS or connection error, ensure `VITE_API_URL` points to where the backend is running (default port **8100**).

## Adding more tools

- **More string tools**: Add handlers under `backend/handlers/`, register routes in `backend/main.go` under `/api/string/`, add entries in `frontend/src/config/sidebarConfig.ts` (with optional `subGroup` for collapsible groups), and wire up UI in `frontend/src/components/StringTools.tsx` and API helpers in `frontend/src/api/stringTools.ts`.
- **Other tool types**: Add new route groups in `main.go` (e.g. `/api/encode/`, `/api/hash/`) and new top-level nav sections and pages in the frontend; same pattern: Go handlers + React that calls the API.
