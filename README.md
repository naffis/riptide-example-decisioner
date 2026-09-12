# riptide-example-decisioner

Public starter for a Riptide **EXTERNAL_ENDPOINT** Decisioner: a small Go HTTP server that accepts
a protojson-shaped `DecideRequest` on `POST /decide` and returns a deterministic floor bump.

Docs: [Extensions](https://riptide.adwave.com/docs/extensions) · Agent paste prompt: [PROMPT.md](./PROMPT.md)

## Run

```sh
go test ./...
go run . -addr 127.0.0.1:8089
```

Smoke:

```sh
curl -sS -X POST http://127.0.0.1:8089/decide \
  -H 'Content-Type: application/json' \
  -d '{
    "request": {
      "point": "DECISION_POINT_FLOOR",
      "tenantId": "018f0000-0000-7000-8000-000000000001",
      "requestId": "sample-request-1",
      "policy": { "floor": "2.0000", "maxBid": "5.0000" }
    }
  }'
```

## Register with Riptide

Deploy the service behind an HTTPS proxy with a valid certificate and a public hostname that
resolves to public IP addresses. Bind to the interface required by your deployment, for example
`go run . -addr :8089`. Replace the example hostname below with your deployed service, then create
a model with `kind: EXTERNAL_ENDPOINT` and config:

```json
{"url":"https://decisioner.example.com/decide"}
```

(`endpoint` is also accepted by serve.) The host rejects plain HTTP, loopback, and private-network
destinations. The localhost command above is a local smoke test. Assign as `SHADOW` before
`BOUNDED_AB` or `LIVE`.

## Notes

- This sample uses plain `encoding/json` with protojson field names so it runs with no private
  monorepo dependency. Production hosts send protojson; unknown fields should be ignored if you
  switch to protobuf libraries later.
- Money is always a decimal string with 4 decimal places.
- Request bodies may be up to 1 MiB; larger bodies receive HTTP 413. Headers, complete request
  reads, and response writes each have a 5-second timeout; idle keep-alive connections expire
  after 30 seconds. These limits suit the starter's quick response; they do not extend the
  host's shorter decision budget.
- On error or timeout the host serves the deterministic baseline Decisioner.
