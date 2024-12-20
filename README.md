# Templ in Cloudflare workers

An example/template project that uses Cloudflare's Workers and Pages service, along with a key-value store to host a templ app.

[SSE in action](https://github.com/user-attachments/assets/ee856356-7b9b-4970-9756-31e8b378b71f)

Features:
  - Counters are updated without a full page reload.
  - Queued requests so all button presses are counted.
  - Worker compiled to web assembly.
  - Uses `htmx` for interactivity, `templ` for page rendering, `tinygo` to compile for WASM and `wrangler` to deploy.
  - SSE support - counter also updates globally without a reload.
