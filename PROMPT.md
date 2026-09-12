# Agent prompt (Decisioner HTTP starter)

Copy this into your coding agent, then append your goal.

```
Clone https://github.com/naffis/riptide-example-decisioner (or use examples/public/decisioner-http in a Riptide checkout).

Build or extend an EXTERNAL_ENDPOINT Riptide Decisioner:
- HTTP POST /decide accepts JSON matching riptide.decision.v1 DecideRequest (protojson field names).
- Return DecideResponse with a price CPM as a decimal string (4 dp), never float64.
- Honor context cancellation; on timeout the host falls back to its deterministic baseline.
- Deploy behind HTTPS with a valid certificate and a hostname resolving to public IP addresses.
- Register with model kind EXTERNAL_ENDPOINT and config {"url":"https://decisioner.example.com/decide"}, replacing the example hostname with your deployment.
- The host rejects plain HTTP, loopback, and private-network destinations; localhost HTTP is for local smoke tests only.
- Keep the 1 MiB request limit and bounded HTTP server timeouts; return HTTP 413 for oversized bodies, including requests without Content-Length.
- Assign as SHADOW first; promote to LIVE only after measured lift.
- Keep secrets out of config; use env vars.

My goal:
(describe the decision point and behavior)
```

Full platform prompt (all extension surfaces): https://riptide.adwave.com/docs/extensions
