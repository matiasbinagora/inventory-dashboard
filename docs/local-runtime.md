# Local web runtime

The web app is intentionally local-only. Start the Go API first (it binds to
`127.0.0.1:8080` by default), then start the Next.js app:

```bash
go run ./cmd/inventory-api
npm install
npm run dev
```

Open `http://127.0.0.1:3000`. Next.js proxies `/api/*` to the Go API at
`127.0.0.1:8080`; the web UI does not duplicate API or business logic. The
dashboard displays only records returned by `GET /api/projects`. No development
seeds, credentials, source paths, private URLs, or sensitive media are bundled.
