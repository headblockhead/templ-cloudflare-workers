# Templ in Cloudflare workers

Uses Cloudflare's Workers and Pages service, along with a key-value store to host a templ app.

Features:
  - Counters are updated without a full page reload
  - Queued requests so all button presses are counted.
  - Worker compiled to web assembly as Go is sadly not supported in Workers
  - Uses `htmx` for interactivity, `templ` for page rendering, `tinygo` to compile for WASM and `wrangler` to deploy.
