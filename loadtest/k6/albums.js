import http from "k6/http";
import { check, sleep } from "k6";

// https://grafana.com/docs/k6/latest/using-k6/http-requests/#group-urls-under-one-tag
const BASE = __ENV.BASE_URL || "http://localhost:8080";

function params(routeName, extra = {}) {
	return {
		...extra,
		// `name` last so callers cannot override the route label (k6 docs: tags.name).
		tags: { ...(extra.tags || {}), name: routeName },
	};
}

function newAlbumID() {
	const jitter = Math.floor(Math.random() * 10_000);
	return String(
		1_000_000_000 + (__VU % 50_000) * 10_000 + (__ITER % 10_000) * 17 + jitter,
	);
}

export const options = {
	vus: Number(__ENV.VUS || 10),
	duration: __ENV.DURATION || "30s",
	thresholds: {
		http_req_failed: ["rate<0.05"],
		http_req_duration: ["p(95)<500"],
	},
};

export default function () {
	const headers = { "Content-Type": "application/json" };

	let res = http.get(`${BASE}/`, params("GET /"));
	check(res, { "GET / 200": (r) => r.status === 200 });

	res = http.get(`${BASE}/albums`, params("GET /albums"));
	check(res, { "GET /albums 200": (r) => r.status === 200 });

	const id = newAlbumID();
	const createBody = JSON.stringify({
		id,
		title: `loadtest-${id}`,
		artist: "k6",
		price: 9.99 + Math.random(),
	});

	res = http.post(`${BASE}/albums`, createBody, params("POST /albums", { headers }));
	check(res, { "POST /albums 201": (r) => r.status === 201 });

	res = http.get(`${BASE}/album/${id}`, params("GET /album/:id"));
	check(res, { "GET /album/:id 200": (r) => r.status === 200 });

	const updateBody = JSON.stringify({
		id,
		title: `loadtest-updated-${id}`,
		artist: "k6",
		price: 19.99,
	});
	res = http.put(`${BASE}/album/${id}`, updateBody, params("PUT /album/:id", { headers }));
	check(res, { "PUT /album/:id 200": (r) => r.status === 200 });

	res = http.del(`${BASE}/album/${id}`, params("DELETE /album/:id", { headers }));
	check(res, { "DELETE /album/:id 200": (r) => r.status === 200 });

	sleep(Number(__ENV.SLEEP_MS || 100) / 1000);
}
