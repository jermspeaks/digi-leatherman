# Digi Leatherman

A full-stack app for string tools (and later other tool types). The **Go backend** exposes an API; the **React + TypeScript** frontend (Vite) calls it.

## Run the app

### Backend (Go API)

From the repo root:

```bash
cd backend
go run .
```

The API listens on **http://localhost:8100**. Endpoints:

- `POST /api/string/url-encode` — body: `{"value": "..."}` → `{"result": "encoded%20string"}`
- `POST /api/string/url-decode` — body: `{"value": "..."}` → `{"result": "decoded string"}`

### Frontend (Vite + React)

In another terminal:

```bash
cd frontend
npm install   # if you haven't already
npm run dev
```

Open **http://localhost:5273**. The UI has “String tools” with **URL encode** and **URL decode**; enter text, click Encode/Decode, and see the result. The frontend uses `VITE_API_URL` (default `http://localhost:8100`); copy `.env.example` to `.env` and change it if your API runs elsewhere.

## Adding more tools

- **More string tools**: Add handlers under `backend/handlers/` (e.g. base64, trim), register routes in `backend/main.go` under `/api/string/`, and add UI in `frontend/src/components/StringTools.tsx` (or new components) and API helpers in `frontend/src/api/stringTools.ts`.
- **Other tool types**: Add new route groups in `main.go` (e.g. `/api/encode/`, `/api/hash/`) and new top-level nav sections and pages in the frontend; same pattern: Go handlers + React that calls the API.
