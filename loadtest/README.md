# Load tests (k6)

Uses [Grafana k6](https://k6.io/) to hit every HTTP route: `GET /`, `GET /albums`, `POST /albums`, `GET /album/:id`, `PUT /album/:id`, `DELETE /album/:id`.

Each request passes **`tags: { name: '<route>' }`**, as in the [k6 HTTP docs](https://grafana.com/docs/k6/latest/using-k6/http-requests/#group-urls-under-one-tag), so dynamic URLs (e.g. `/album/123`) still aggregate under one logical name in metrics and dashboards.

For **per-route tables in the terminal**, use e.g. `k6 run --summary-export=summary.json ...` and inspect the export, or send output to **k6 Cloud** / **Grafana** / **JSON** (`--out json=results.json`) and filter by the `name` tag.

POST uses a **random numeric string** `id` so it works with in-memory storage and typical **MySQL `INT`** primary keys. You can delete those rows yourself afterward.

## Install k6

- **Windows (Chocolatey):** `choco install k6`
- **macOS:** `brew install k6`
- **Linux / other:** https://grafana.com/docs/k6/latest/set-up/install-k6/

## Run (server must already be listening)

From repo root:

```bash
cd loadtest/k6
k6 run -e BASE_URL=http://localhost:8080 albums.js
```

Optional environment:

| Variable    | Default   | Description                |
|------------|-----------|----------------------------|
| `BASE_URL` | `http://localhost:8080` | API origin         |
| `VUS`      | `10`      | Virtual users              |
| `DURATION` | `30s`     | Test length                |
| `SLEEP_MS` | `100`     | Pause between iterations (ms) |

Example:

```bash
k6 run -e BASE_URL=http://127.0.0.1:8080 -e VUS=25 -e DURATION=1m albums.js
```

## Docker (no local k6 install)

From `loadtest/k6` (adjust `BASE_URL` if the app runs on the host: use `host.docker.internal` on Docker Desktop instead of `localhost`):

```bash
docker run --rm -i -v "${PWD}:/scripts" grafana/k6 run -e BASE_URL=http://host.docker.internal:8080 /scripts/albums.js
```

On Linux you may need `--add-host=host.docker.internal:host-gateway` instead.
