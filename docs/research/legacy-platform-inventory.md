# Legacy platform inventory

Research dates: 2026-10-24 (source), 2026-10-26 (production, see "Production inventory"). Scope: `bb_api/` source plus read-only production inspection. Jira holds the canonical decision; this file is supporting research.

## Evidence boundaries and access

- Main API advertises `/api/bb/*.php`; the frontend is `bb_api/app.js`, and `index.php` lists public/user/admin contracts ([bb_api/index.php:8](../../bb_api/index.php#L8)). Actual web-server mappings, deployed revision, host ownership, enabled endpoints and domains require operator confirmation.
- Database selection is `BB_DB_DRIVER` / `BB_DB_DSN`; PostgreSQL additionally uses `BB_DB_HOST`, `BB_DB_PORT`, `BB_DB_NAME`, `BB_DB_USER`, `BB_DB_PASSWORD`, `BB_DB_SSLMODE`. Default is SQLite at `MKF_STORAGE_PATH/bb.sqlite`; MySQL is recognized but rejected by the connector ([helpers.php:1743](../../bb_api/lib/helpers.php#L1743)). Source does not establish the production driver.
- Environment lookup supports `BB_ENV_FILE`, `/etc/machkinfach/.env`, `/var/www/.machkinfach.env` and repository-adjacent files ([helpers.php:19](../../bb_api/lib/helpers.php#L19)). These are configuration lookup paths, not proof that any exists or grants access. No secret values are included here.
- Do not use `bb_db()` for read-only inventory: it performs schema, encryption and session migrations and automatic OCR retention deletion ([helpers.php:1798](../../bb_api/lib/helpers.php#L1798), [helpers.php:1812](../../bb_api/lib/helpers.php#L1812)). Even `bb_admin_db_readonly()` calls `bb_db()` before applying read-only mode ([admin/common.php:13](../../bb_api/admin/common.php#L13)). Use a direct raw read-only database connection or snapshot instead.

## Readers, writers and logical ownership

| Plane | Readers / writers visible in source | Ownership and relationship |
|---|---|---|
| Main PHP SQL database | Auth/register/login; analysis and WhatsApp ingest; Telegram BP bridge; BP rewards/grants/purchases/Stripe webhook; admin; expert/partner portal; account and analysis deletion | Main `users.id` relates to sessions, letters and many `user_id` columns. Referrals have inviter/invitee owners; experts add partner ownership/claim and user consent. Not every logical relationship has a declared FK. Source schema inventory below records tables/columns. |
| Main filesystem | Helpers, analysis runner, PDF generators, solutions/partners administration, import, download links and deletion endpoints | User JSON and per-user letter folders coexist with SQL records. Source originals and generated PDFs are referenced by path; copying SQL alone is insufficient. |
| Telegram DataStore | `node.js` reads/writes `users`, `stats`, `rating`, `broadcast_job`; optional referral and broadcast workers may also write | Bot uses `chatId`, plan/subscription, locale, activity and sender profiles in its own user array. Main BP state is shared via signed PHP bridge. Must reconcile bot-only users with main accounts, rather than assume all are paid customers or already imported. |
| Optional Telegram runtime | Injected `globalThis.__bbRuntime` handles enqueue/stats/job/list/retry | Durable job backend and runtime implementation are absent from this folder. Existence, owners, persistence, deployment and retry schedule are unverified. |
| Buddy Manager | Separate PHP API, resident/letter/user/profile/form flows | Own `data/app.sqlite`; `app_users` identity is independent of main `users`. Resident-linked letters do not carry a direct main account FK. Source auto-creates an embedded admin identity; real tenant/ownership boundary must be confirmed. |

Evidence: SQL ownership [helpers.php:1851](../../bb_api/lib/helpers.php#L1851); expert claim/consent columns [expert_requests.php:14](../../bb_api/lib/expert_requests.php#L14); Telegram stores [node.js:73](../../bb_api/node.js#L73), bot user shape [node.js:777](../../bb_api/node.js#L777), signed shared BP call [node.js:208](../../bb_api/node.js#L208), bridge upsert [tg_bp_bridge.php:37](../../bb_api/tg_bp_bridge.php#L37); Buddy Manager DB [db.php:8](../../bb_api/buddy-manager/lib/db.php#L8), embedded identity [auth.php:24](../../bb_api/buddy-manager/lib/auth.php#L24).

## Stored objects beyond SQL

| Location (source defaults, relative to storage unless noted) | Contents and evidence |
|---|---|
| `MKF_STORAGE_PATH` or `bb_api/bb_storage` | Main root; `logs/`, `users/`, `letters/`, `analysis_requests/` ([helpers.php:280](../../bb_api/lib/helpers.php#L280)). |
| `users/<id>.json` | User storage JSON; atomic temporary write + rename ([helpers.php:1578](../../bb_api/lib/helpers.php#L1578), [helpers.php:1593](../../bb_api/lib/helpers.php#L1593)). Admin also supports `users/*/profile.json`, so do not assume one historical layout ([common.php:52](../../bb_api/admin/common.php#L52)). |
| `letters/<user-id>/<timestamp>.json` | Analysis metadata snapshot stored alongside SQL ([helpers.php:1723](../../bb_api/lib/helpers.php#L1723), [analysis_runner.php:657](../../bb_api/lib/analysis_runner.php#L657)). |
| `analysis_requests/source/<request-id>_source.<ext>` | Retained uploaded PDF/image/text originals, SQL `source_file_path` references ([analysis_runner.php:108](../../bb_api/lib/analysis_runner.php#L108)). |
| Generated PDFs / temporary export files | `pdf.php` stores final PDF in letter directory and updates SQL path; other PDF generators create temporary HTML/PDF files ([pdf.php:326](../../bb_api/pdf.php#L326)). Public export files additionally live at the web root `storage/solutions_downloads`, outside the main storage root. Public PDF link endpoints also clean files older than one day ([paper_letter_download_link.php:32](../../bb_api/paper_letter_download_link.php#L32), [solutions_download_link.php:35](../../bb_api/solutions_download_link.php#L35)). |
| `solutions_library/catalog.json`, `solutions_library/files/` | Solution catalog, PDF files, download counters ([solutions_storage.php:5](../../bb_api/lib/solutions_storage.php#L5)). |
| `partners/partners.json`, `partners/photos/` | Partner metadata, specialist access codes and images ([partners_storage.php:5](../../bb_api/lib/partners_storage.php#L5), [partners_storage.php:36](../../bb_api/lib/partners_storage.php#L36)). Treat catalog as sensitive configuration. |
| `admin/analysis_counter.json` | Counter presentation/configuration ([analysis_counter.php:15](../../bb_api/lib/analysis_counter.php#L15)). |
| Legacy import output | `full_datastore_stats.json`, `full_datastore_rating.json`; import is an additional potential writer ([import_full_datastore_export.php:255](../../bb_api/import_full_datastore_export.php#L255)). |
| `buddy-manager/data/app.sqlite`, `buddy-manager/uploads/` | Separate durable DB and uploaded/generated document files ([db.php:8](../../bb_api/buddy-manager/lib/db.php#L8), [upload.php:21](../../bb_api/buddy-manager/api/upload.php#L21), [pdf.php:67](../../bb_api/buddy-manager/api/pdf.php#L67)). PHP session storage is another runtime dependency ([auth.php:4](../../bb_api/buddy-manager/lib/auth.php#L4)). |

PII, OCR, analysis JSON, blind indexes and session hashes rely on `BB_USER_PII_KEY`, `BB_PARSED_TEXT_KEY`, `BB_SESSION_TOKEN_PEPPER`; key continuity must accompany data preservation. Plaintext and encrypted columns coexist in source; production encryption coverage is unknown ([helpers.php:310](../../bb_api/lib/helpers.php#L310), [helpers.php:316](../../bb_api/lib/helpers.php#L316), [helpers.php:591](../../bb_api/lib/helpers.php#L591)).

## Imported DataStore reconciliation

The importer accepts an export of key/value records, reads `users`, `stats`, `rating`, historical `last_parsed_<chat-id>` keys and user-embedded `parsedText` ([import_full_datastore_export.php:30](../../bb_api/import_full_datastore_export.php#L30)). This historical export shape differs from the current bot's in-memory OCR design; preserve and inventory export files separately rather than infer they are obsolete.

For users present in the export user list, it looks up existing main accounts by Telegram blind index or plaintext identifier, then updates/creates accounts. For text-only Telegram identifiers absent from that list, it directly creates a new account without that same existing-account lookup. Letter deduplication is by `(user_id, analysis_text)`, and can restore encrypted OCR into an existing matching row ([import_full_datastore_export.php:91](../../bb_api/import_full_datastore_export.php#L91), [import_full_datastore_export.php:143](../../bb_api/import_full_datastore_export.php#L143), [import_full_datastore_export.php:163](../../bb_api/import_full_datastore_export.php#L163), [import_full_datastore_export.php:196](../../bb_api/import_full_datastore_export.php#L196)). Therefore importer existence does not prove an import ran, is complete, preserves every bot field, or is globally idempotent. Reconcile bot/export/main identities, duplicate users, source text and balances before relying on imported rows as the sole copy. Never replay this tool as a read-only inventory probe.

The tool calls `bb_db()` before disabling retention, so initialization can already perform retention/migration writes. Its source-default input path is a host filesystem hint, not verified access ([import_full_datastore_export.php:11](../../bb_api/import_full_datastore_export.php#L11), [import_full_datastore_export.php:73](../../bb_api/import_full_datastore_export.php#L73)).

## Jobs, integrations and in-flight operations

- Web/WhatsApp analyses create `analysis_requests`, retain source files and call the analysis executor synchronously; failures are persisted, and admin retry creates/uses request lineage ([analyze.php:103](../../bb_api/analyze.php#L103), [analyze.php:125](../../bb_api/analyze.php#L125), [admin/analysis_retry.php:1](../../bb_api/admin/analysis_retry.php#L1)). Live pending/running/failed counts, active HTTP work and retry ownership are unknown.
- OCR uses Google Vision and language generation uses OpenAI chat completions ([lib/openai.php:28](../../bb_api/lib/openai.php#L28), [lib/openai.php:102](../../bb_api/lib/openai.php#L102)). Buddy Manager has its own integration wrappers.
- Telegram source is a Pipedream-style component receiving `steps.trigger.event.body`. It supports an optional injected OCR queue with enqueue/job/list/retry/stat methods; no scheduler or queue implementation is included ([node.js:6](../../bb_api/node.js#L6), [node.js:1613](../../bb_api/node.js#L1613), [node.js:2571](../../bb_api/node.js#L2571)). Telegram webhook registration and actual runtime URL must be inspected externally.
- Bot broadcast writes a recipient snapshot, index and sent/failed/blocked counters to `broadcast_job`, then POSTs `BROADCAST_WORKER_URL`; `/start` can invoke `REFERRAL_WORKFLOW_URL`. These workers are external, so enumerate them and their datastore access before switching writers ([node.js:1905](../../bb_api/node.js#L1905), [node.js:2444](../../bb_api/node.js#L2444)).
- Bot BuddyPoints bridge authenticates using timestamped HMAC and can upsert accounts and write balances/analysis state; both bot and web therefore participate in shared user state ([tg_bp_bridge.php:7](../../bb_api/tg_bp_bridge.php#L7), [tg_bp_bridge.php:37](../../bb_api/tg_bp_bridge.php#L37)).
- Stripe webhook accepts checkout completion and async payment success, updates purchase/balance/ledger/message records. Preserve outstanding checkout client references and reconcile pending payments and provider event replay with the operator ([stripe_webhook.php:34](../../bb_api/stripe_webhook.php#L34), [stripe_webhook.php:100](../../bb_api/stripe_webhook.php#L100)).
- Pingen uploads files and creates letters with `auto_send=true`, records `pingen_id` and local `sent`, and refunds BP on reported failures. Source `sent` means API create succeeded, not physical delivery confirmed. A timeout after an external send can require provider-side reconciliation ([paper_letter_send.php:202](../../bb_api/paper_letter_send.php#L202), [paper_letter_send.php:255](../../bb_api/paper_letter_send.php#L255)). Buddy Manager separately sends through Pingen and stores provider references ([send.php:25](../../bb_api/buddy-manager/api/send.php#L25)).
- Email uses SMTP or PHP mail/sendmail fallback; Buddy Manager has its own SMTP sender. WhatsApp code delivery uses configurable bridge `/send-message`; ingest accepts internal authenticated multipart input ([helpers.php:167](../../bb_api/lib/helpers.php#L167), [buddy-manager/lib/mail.php:42](../../bb_api/buddy-manager/lib/mail.php#L42), [auth_request_whatsapp_code.php:63](../../bb_api/auth_request_whatsapp_code.php#L63), [whatsapp_ingest.php:587](../../bb_api/whatsapp_ingest.php#L587)). Buddy Manager fax uses `BB_FAX_API_URL` / token ([fax.php:24](../../bb_api/buddy-manager/lib/fax.php#L24)). Provider accounts, pending messages/faxes and WhatsApp bridge state are unverified.
- Parsed-text retention is triggered on normal DB initialization, in addition to manual admin purge. No cron/systemd/deployment schedule definitions were found under `bb_api`; absence in this snapshot does not establish absence in production ([helpers.php:827](../../bb_api/lib/helpers.php#L827), [helpers.php:1812](../../bb_api/lib/helpers.php#L1812)).
- Manual writer tools also include SQLite→PostgreSQL migration, datastore import, sender reprocessing and backfill. They must be accounted for in writer control, even if not run by this investigation (see endpoint/source manifest).

## Telegram flows and command surface

User paths include start/referral, language selection (UA/RU/DE/EN), document/photo intake, queued or inline OCR, AI letter explanation, follow-up questions, response template generation, rating/feedback, support, invite and restart. Parsed text lives in `sessionMemory` with TTL rather than durable user records; sender context and profile state do persist in bot users ([node.js:3](../../bb_api/node.js#L3), [node.js:51](../../bb_api/node.js#L51), [node.js:777](../../bb_api/node.js#L777), [node.js:1684](../../bb_api/node.js#L1684), [node.js:2550](../../bb_api/node.js#L2550)).

Commands visible in dispatch: `/start`, `/stats`, `/queue`, `/job`, `/jobs`, `/retry`, `/activity`, `/setplan`, `/broadcast`, `/audience`, `/reply`. Admin commands check configured/hardcoded administrator identity and optionally a secret. Do not invoke mutating commands for inventory. Callback flows include language codes, `ask_question`, `start_over`, `leave_feedback`, `rate_1`…`rate_5` ([node.js:1480](../../bb_api/node.js#L1480), [node.js:2049](../../bb_api/node.js#L2049), [node.js:2311](../../bb_api/node.js#L2311), [node.js:2403](../../bb_api/node.js#L2403), [node.js:2685](../../bb_api/node.js#L2685)). This is source dispatch coverage, not live Telegram menu configuration.

## Production inventory (2026-10-26)

Read-only inspection of host `srv1433675` (187.124.12.51) as OS user `rene`; DB queried as `www-data` (peer auth, via a temporary `sudoers.d/rene-psql-readonly` rule) inside `default_transaction_read_only` transactions. Aggregate counts only; no customer rows or secret values were read out.

### Host and runtime

- Ubuntu 24.04, 8 GB RAM, 96 GB disk (21 % used). Native nginx + php8.3-fpm + PostgreSQL 16.15; legacy app is not containerised. The same host also runs a Coolify control plane in Docker.
- Deployed `/var/www/mk/bb_api` is byte-identical to repo `bb_api/` for all 155 PHP/JS files (repo-only: `app.js`, `node.js`). The source inventory above therefore describes deployed code.
- Config: `/var/www/mk/.env` and `/var/www/.machkinfach.env` hold the same non-secret settings: `BB_DB_DRIVER=pgsql`, DB `machkinfach` over `/var/run/postgresql`, `BB_DB_USER=www-data` (no password; peer auth), `MKF_STORAGE_PATH=/var/www/mk/bb_storage/bescheidbuddy`, retention 30 days, `OPENAI_MODEL=gpt-4o-mini`, SMTP `smtp.hostinger.com` from `kontakt@machkinfach.de`. Configured providers: OpenAI, Google Vision, Stripe, Pingen, PerfectCorp, Telegram, WhatsApp bridge.
- nginx `machkinfach.de`: `/` → `/var/www/mk/root`, `/app/` → `bb_app`, `/api/bb/` → `bb_api` (GET/POST/OPTIONS; import/backfill scripts denied), `/storage/bb/` → `bb_storage` (PHP/JSON/DB extensions denied, other files served). Also `/mb/` (`/srv/heim_buddy_mvp`), `/photography`, `/g`. Separate vhosts `mb.machkinfach.de` and `test.machkinfach.de` (default server) serve `heim_buddy_mvp` over HTTP.
- systemd: `bb-whatsapp-bridge.service` (`/var/www/mk/whatsapp-bridge`, 96 MB) is live. `tg-bot.service` (`node /var/opt/tg_bot/runtime.mjs`, user `tgbot`, directory unreadable) is running and is the **live bot** (polling); see "Live Telegram runtime". Pipedream holds a stale DataStore.
- No app cron or timer visible to `rene`; root and `www-data` crontabs need sudo. No backup job found.

### Database `machkinfach` (34 MB, 27 tables)

| Table | Rows | Note |
|---|---|---|
| users | 1,733 | 32-hex ids; created 2025-11-04 → 2026-10-23 |
| letters | 773 | owned by 332 users; OCR and analysis 100 % encrypted (`*_enc`), 0 plaintext |
| sender_stats | 1,298 | |
| security_audit | 2,322 | |
| bp_ledger | 820 | 255 users with BP activity; 218 users hold earned BP (Σ 3,906) |
| user_messages | 625 | |
| sessions | 419 | 4 unexpired; not migrated |
| analysis_requests | 130 | 111 success, 19 failed, none queued/running |
| bp_purchase_requests | 31 | 3 paid (2 manual, 1 Stripe) → **2 paying customers**; 30 pending (Stripe 26, manual 4), all from 2026-04/05, i.e. stale |
| paper_letter_deliveries | 21 | 15 sent, 6 failed, none pending |
| support_messages | 18 | 15 new (unanswered) |
| partner_expert_requests | 4 | all `new`, consent `none` |
| bp_grants / referrals / social_review_rewards | 10 / 2 / 2 | |
| auth_* / admin_sessions / auth_attempts | — | transient, not migrated |
| usage_monthly, social_links_config, parsed_text_retention_settings | 0 | |

- **Identity channels:** 39 email (all encrypted), 17 with a password, 254 Telegram (all encrypted), 4 WhatsApp. Zero blind-index duplicates for email, Telegram or WhatsApp. Zero orphans across all `user_id` references and `analysis_requests.letter_id`.
- **1,436 users have no identity in SQL** (no email/Telegram/WhatsApp). Created 2025-11 (1,058), 2025-12 (248), 2026-01…03 (130): the DataStore import. 257 own letters; 2 have `purchase_paid` ledger entries. Their Telegram link exists only as `chatId` in `users/<id>/profile.json`.
- File references: 20 `letters.pdf_path` absolute under `/var/...`; 172 letters and 130 analysis requests reference `source_file_path` relative to storage.

### Storage `bb_storage/bescheidbuddy` (414 MB, 2,500 files)

- `analysis_requests/` 220 MB / 126 files; `letters/` 155 MB / 568 files (492 JSON, 54 JPG, 22 PDF); `users/` 31 MB.
- `users/`: 1,435 `<id>/profile.json` (keys: `userId`, `chatId`, `plan`, `subscriptionUntil`, `used`/`limit`, `questionLimit`/`usedQuestions`, `lang`, `firstSeen`/`lastSeen`, **`parsedText` in 1,363**, `lastOcrAt`) plus 322 `<id>.json` in the newer shape (`email`, `tg_user_id`, `wa_*`, `sessions`), plus a few `uploads/`, `letters.json`, `usage-*.json`.
- `parsedText` in profile JSON is plaintext OCR outside SQL encryption and outside the 30-day SQL retention purge.
- `legacy/` holds the 2026-03-08 DataStore import metadata (`full_datastore_stats.json`, `full_datastore_rating.json`, `stats_extracted_once_*.json`).
- Leftover `bb.sqlite` (2.7 MB, 2026-03-08), pre-PostgreSQL copy.
- Sibling non-BescheidBuddy trees in `bb_storage`: `app`, `galleries`, `hair_tryon_spaces`, `orders`, `photography-cache`.

### Other data planes

- Buddy Manager: `bb_api/buddy-manager/data/app.sqlite` (128 KB, 2026-05-22). Separate `heim_buddy_mvp` app: `/srv/heim_buddy_mvp/.../backend/data/app.sqlite` (144 KB, 2026-06-02), 18 MB total.
- `/var/www/bb` (14 MB, 928 files): purpose unidentified.

### Pipedream DataStore (2026-10-26)

Read via a temporary HTTP workflow in project `proj_rKsbdNG` returning aggregates only.

- 75 keys: `users`, `stats` (141 top-level keys), `rating` (count 9), and 72 `last_parsed_<chatId>`. No `broadcast_job` key (no broadcast pending or ever persisted).
- `users`: array of 1,441, no duplicate or missing `chatId`. Plans: `free` 1,213, `marrychristmas` 135, `buddys` 3, none 90; **0 active subscriptions**. Languages: ru 1,088, ua 228, de 39, en 2, none 84. 259 users carry non-empty `parsedText`.
- `firstSeen`: 2025-11 1,158, 2025-12 174, 2026-01…04 21, none 88. **Latest `lastSeen` is 2026-04**: no user activity has been written to this DataStore since April 2026.
- 1,441 DataStore users vs 1,435 imported `profile.json` (import 2026-03-08): roughly 6 bot users were never imported into the main platform.

### Live Telegram runtime (2026-10-26)

- `getWebhookInfo` returns an empty URL: the bot **polls**. The legacy `.env` `TELEGRAM_BOT_TOKEN` is rejected by Telegram (401), so PHP holds a stale token.
- A test text message left no trace in PHP, PostgreSQL or the Pipedream DataStore. The Pipedream DataStore stayed unchanged (no `lastSeen` after 2026-04): **Pipedream is a stale snapshot, not the live bot.**
- A test document upload produced, within 17 s: a `bp_ledger` `spend_analyze` row (13:54:30 UTC), an `analysis_ok` log entry with `source: telegram`, a new `letters` row, an updated `users/<id>.json` and a new `letters/<id>/<timestamp>.json`. No `analysis_requests` row (latest is 2026-09-09): the Telegram path bypasses request tracking.
- So the live bot is `tg-bot.service` (`/var/opt/tg_bot/runtime.mjs`, polling), and its durable outputs land in the main PHP platform. Whatever it keeps locally in `/var/opt/tg_bot` (user list, language, plan, session state) is still unread (needs sudo).

### Live bot state: `/var/opt/tg_bot` (2026-10-26)

- Code: `node.js` is byte-identical to repo `bb_api/node.js`; it is wrapped by `runtime.mjs` (2026-07-09, not in the repo), which provides polling, the SQLite-backed store and the OCR job queue (the `__bbRuntime` absent from source). `.env` keys: `TELEGRAM_TOKEN`, `OPENAI_KEY`, `GOOGLE_VISION_KEY`, `OCR_WORKERS`.
- Data: `data/bot.sqlite` (1.8 MB) in WAL mode with an un-checkpointed 3.9 MB `-wal`; any copy must include `-wal`/`-shm` or be taken via the SQLite backup API.
- `jobs` (OCR queue): 512 rows since 2026-04-12 (491 done, 21 failed, 0 unfinished), 106 distinct chats. Stores full raw Telegram updates (personal data). 2026-04-12 is the Pipedream → host cutover date.
- `kv`: `__telegram_offset` (polling offset), `users` (**1,941**, no duplicate `chatId`; plans free 1,821 / marrychristmas 118 / buddys 2; ru 1,526, ua 351, de 58, en 6; 436 first seen in 2026-04; 67 carry the shared-BP fields, 107 sender profiles; no `parsedText` field), `stats` (282 keys), `rating`, `broadcast_job` (last run 2026-07-09: 1,927 recipients, 1,236 sent, **691 blocked**, status `done`).

### Identity planes compared

| Store | Telegram users | Status |
|---|---|---|
| `bot.sqlite` `kv.users` | 1,941 | live, authoritative bot user list |
| Pipedream DataStore `users` | 1,441 | frozen since 2026-04 |
| Main `profile.json` (imported) | 1,435 | import 2026-03-08; link only via `chatId` |
| Main SQL `users.tg_user_id_enc` | 254 | web/Telegram login and bridge |

Roughly 500 Telegram users exist only in `bot.sqlite`. The overlap between the 254 SQL Telegram ids and the 1,435 imported profiles is unknown without decryption.

### Schedules and backups

- No crontab for root, `www-data`, `tgbot` or `postgres`; `/etc/cron.*` contain only OS/certbot/docker housekeeping. Parsed-text retention runs only via `bb_db()` on requests.
- **No backups exist**: `/var/backups` is empty of app data, and no dump job, backup directory or schedule was found for PostgreSQL, `bb_storage` or `bot.sqlite`.

### Not inventoried

- Contents of the MachKinfach side apps on the same host (`heim_buddy_mvp`, `/photography`, `/g`, hair try-on, orders, `/var/www/bb`) beyond size and routing.
- Provider-side state (Stripe dashboard, Pingen physical delivery status, WhatsApp bridge session).

## Source schema inventory

These are source-declared tables with CREATE columns and additive columns; the two separately scoped `letters` tables belong to different databases. Main `users` BP/referral/welcome columns are added later in the same helper migration rather than its CREATE block. This is a column-name inventory, with types, constraints and index definitions retained in the linked source; actual production schema and constraints need direct metadata verification. SQL statements live inline rather than in a standalone migration directory. Main indexes are declared at [helpers.php:2336](../../bb_api/lib/helpers.php#L2336).

- `users` — `id, email, email_enc, email_bidx, pass_hash, tg_user_id, tg_user_id_enc, tg_user_id_bidx, tg_username, tg_username_enc, tg_username_bidx, wa_chat_id, wa_chat_id_enc, wa_chat_id_bidx, wa_name, wa_name_enc, wa_name_bidx, wa_phone, wa_phone_enc, wa_phone_bidx, signup_bp_rewarded_at, first_analysis_bp_rewarded_at, created_at, monthly_bp_total, monthly_bp_remaining, monthly_bp_reset_at, earned_bp, referral_code, referred_by_user_id, social_review_rewarded_at, welcome_message_sent_at`. [bb_api/lib/helpers.php:1851](../../bb_api/lib/helpers.php#L1851).
- `sessions` — `token, token_hash, user_id, created_at, expires_at, ip, ua, last_seen`. [bb_api/lib/helpers.php:1877](../../bb_api/lib/helpers.php#L1877).
- `usage_monthly` — `user_id, ym, analyses_used, questions_used`. [bb_api/lib/helpers.php:1888](../../bb_api/lib/helpers.php#L1888).
- `letters` — `id, user_id, created_at, sender_name, sender_address, sender_region, sender_category, sender_subcategory, sender_department, sender_contact_person, sender_street, sender_house_number, sender_postal_code, sender_city, sender_state, sender_country, sender_phone, sender_email, sender_website, sender_reference_number, sender_tax_or_case_number, normalized_sender_name, normalized_sender_key, sender_confidence, address_confidence, category_confidence, topic_category, topic_title, urgency_level, filename, mime, size, ocr_text, ocr_text_enc, ocr_text_deleted_at, analysis_json, analysis_json_enc, analysis_text, pdf_path, source_file_path, source_file_mime`. [bb_api/lib/helpers.php:1896](../../bb_api/lib/helpers.php#L1896).
- `analysis_requests` — `id, user_id, parent_request_id, requested_by_admin_id, request_origin, status, source, lang, pages_meta, model, filename, mime, size, source_file_path, source_file_mime, ocr_text_enc, ocr_text_deleted_at, letter_id, error_code, error_message, error_http, error_debug_id, raw_output_preview, retry_used, bp_charged, created_at, updated_at, completed_at`. [bb_api/lib/helpers.php:1939](../../bb_api/lib/helpers.php#L1939).
- `parsed_text_retention_settings` — `id, auto_delete_enabled, retention_days, updated_at, updated_by`. [bb_api/lib/helpers.php:1972](../../bb_api/lib/helpers.php#L1972).
- `sender_stats` — `id, user_id, sender_name, sender_address, sender_region, sender_category, sender_subcategory, sender_city, normalized_sender_name, normalized_sender_key, counter, last_seen`. [bb_api/lib/helpers.php:1980](../../bb_api/lib/helpers.php#L1980).
- `auth_attempts` — `k, cnt, first_ts, last_ts`. [bb_api/lib/helpers.php:1996](../../bb_api/lib/helpers.php#L1996).
- `auth_email_codes` — `id, user_id, email_enc, email_bidx, code_hash, attempts, max_attempts, expires_at, used_at, created_at, ip, ua`. [bb_api/lib/helpers.php:2003](../../bb_api/lib/helpers.php#L2003).
- `auth_whatsapp_codes` — `id, user_id, wa_chat_id_enc, wa_chat_id_bidx, route_chat_id_enc, route_chat_id_bidx, code_hash, attempts, max_attempts, expires_at, used_at, created_at, ip, ua`. [bb_api/lib/helpers.php:2018](../../bb_api/lib/helpers.php#L2018).
- `auth_magic_links` — `id, user_id, token_hash, purpose, expires_at, used_at, created_at, ip, ua`. [bb_api/lib/helpers.php:2035](../../bb_api/lib/helpers.php#L2035).
- `bp_ledger` — `id, user_id, event_type, amount, monthly_delta, earned_delta, meta_json, created_at`. [bb_api/lib/helpers.php:2048](../../bb_api/lib/helpers.php#L2048).
- `referrals` — `id, inviter_user_id, invitee_user_id, referral_code, created_at, signup_rewarded_at, analysis_rewarded_at`. [bb_api/lib/helpers.php:2059](../../bb_api/lib/helpers.php#L2059).
- `social_review_rewards` — `user_id, platform, review_link, meta_json, created_at`. [bb_api/lib/helpers.php:2069](../../bb_api/lib/helpers.php#L2069).
- `bp_grants` — `id, admin_id, user_id, amount, reason, created_at`. [bb_api/lib/helpers.php:2077](../../bb_api/lib/helpers.php#L2077).
- `bp_purchase_requests` — `id, user_id, package_id, bp_amount, status, note, payment_provider, stripe_checkout_session_id, stripe_event_id, paid_at, created_at, updated_at`. [bb_api/lib/helpers.php:2086](../../bb_api/lib/helpers.php#L2086).
- `user_messages` — `id, user_id, title, body, link, is_read, from_admin, message_type, created_at, read_at`. [bb_api/lib/helpers.php:2101](../../bb_api/lib/helpers.php#L2101).
- `social_links_config` — `id, instagram_url, tiktok_url, google_reviews_url, facebook_url, youtube_url, x_url, updated_at`. [bb_api/lib/helpers.php:2114](../../bb_api/lib/helpers.php#L2114).
- `welcome_message_settings` — `id, enabled, title, body, content_i18n_json, link, updated_at, updated_by`. [bb_api/lib/helpers.php:2125](../../bb_api/lib/helpers.php#L2125).
- `support_messages` — `id, user_id, title, body, status, created_at, updated_at`. [bb_api/lib/helpers.php:2136](../../bb_api/lib/helpers.php#L2136).
- `partner_clicks` — `id, user_id, partner_id, category, subcategory, state, partner_name, target_url, lang, source_page, created_at`. [bb_api/lib/helpers.php:2146](../../bb_api/lib/helpers.php#L2146).
- `paper_letter_deliveries` — `id, user_id, pingen_id, status, error_code, error_message, delivery_type, delivery_product, address_position, sender_name, recipient_name, recipient_address, file_name, file_size, page_count, bp_cost, created_at, updated_at`. [bb_api/lib/helpers.php:2160](../../bb_api/lib/helpers.php#L2160).
- `admin_sessions` — `id, token_hash, actor_id, created_at, expires_at, ip, ua, last_seen`. [bb_api/lib/helpers.php:2181](../../bb_api/lib/helpers.php#L2181).
- `security_audit` — `id, event_type, severity, actor_type, actor_id, ip, user_agent, payload_json, created_at`. [bb_api/lib/helpers.php:2192](../../bb_api/lib/helpers.php#L2192).
- `security_alerts` — `id, rule_name, actor_key, severity, status, message, payload_json, created_at, last_seen_at, trigger_count`. [bb_api/lib/helpers.php:2204](../../bb_api/lib/helpers.php#L2204).
- `partner_expert_requests` — `id, support_id, user_id, letter_id, partner_id, partner_name, lang, topic, organization, organization_city, organization_region, deadline, summary_text, status, created_at, updated_at, specialist_note, organization_region_key, routing_scope, published_for_experts, claimed_by_partner_id, claimed_by_partner_name, claimed_at, user_offer_sent_at, user_consent_status, user_consent_at, user_contact_name, user_contact_phone, user_contact_email, user_contact_telegram, secure_docs_requested, secure_docs_requested_at, secure_doc_token, secure_doc_token_created_at, secure_doc_token_used_at, expert_user_rating, expert_user_comment, expert_user_rated_at`. [bb_api/lib/expert_requests.php:15](../../bb_api/lib/expert_requests.php#L15).
- `partner_portal_presence` — `partner_id, last_seen, last_login_at`. [bb_api/partner_portal.php:27](../../bb_api/partner_portal.php#L27).
- `bewoners` — `id, name, first_name, last_name, organization, address, phone, fax, room, contact_name, contact_address, created_at`. [bb_api/buddy-manager/lib/db.php:13](../../bb_api/buddy-manager/lib/db.php#L13).
- `letters` — `id, type, bewoner_id, sender, recipient_address, recipient_email, recipient_fax, topic, subcategory, subject, context, body, signature, provider_status, provider_reference, file_path, created_at, updated_at`. [bb_api/buddy-manager/lib/db.php:28](../../bb_api/buddy-manager/lib/db.php#L28).
- `app_users` — `id, username, password_hash, role, created_at, updated_at`. [bb_api/buddy-manager/lib/db.php:50](../../bb_api/buddy-manager/lib/db.php#L50).
- `sender_profiles` — `id, user_id, first_name, last_name, organization, address, phone, email, fax, created_at, updated_at`. [bb_api/buddy-manager/lib/db.php:59](../../bb_api/buddy-manager/lib/db.php#L59).
- `form_help_sessions` — `id, token, user_id, original_name, file_path, mime, extracted_text, analysis_json, created_at, updated_at`. [bb_api/buddy-manager/lib/db.php:74](../../bb_api/buddy-manager/lib/db.php#L74).

## Complete PHP entry-point and maintenance source manifest

This manifest lists every PHP file outside library directories, grouped by surface; it deliberately includes bootstrap/common files and maintenance scripts because their hosting exposure is unknown. File existence does not prove live HTTP exposure. The advertised contract is [index.php:8](../../bb_api/index.php#L8); this manifest also catches source entry points omitted there. A GET/list-looking endpoint may still initialize/migrate SQL, update access timestamps or remove expired files.

### Main

- [account_delete.php](../../bb_api/account_delete.php#L1)
- [admin_bp_grant.php](../../bb_api/admin_bp_grant.php#L1)
- [admin_message_send.php](../../bb_api/admin_message_send.php#L1)
- [admin_users_search.php](../../bb_api/admin_users_search.php#L1)
- [analysis_counter_public.php](../../bb_api/analysis_counter_public.php#L1)
- [analysis_delete.php](../../bb_api/analysis_delete.php#L1)
- [analysis_get.php](../../bb_api/analysis_get.php#L1)
- [analysis_list.php](../../bb_api/analysis_list.php#L1)
- [analysis_report_error.php](../../bb_api/analysis_report_error.php#L1)
- [analysis_source_file.php](../../bb_api/analysis_source_file.php#L1)
- [analyze.php](../../bb_api/analyze.php#L1)
- [auth_magic_login.php](../../bb_api/auth_magic_login.php#L1)
- [auth_request_code.php](../../bb_api/auth_request_code.php#L1)
- [auth_request_whatsapp_code.php](../../bb_api/auth_request_whatsapp_code.php#L1)
- [auth_telegram.php](../../bb_api/auth_telegram.php#L1)
- [auth_verify_code.php](../../bb_api/auth_verify_code.php#L1)
- [auth_verify_whatsapp_code.php](../../bb_api/auth_verify_whatsapp_code.php#L1)
- [backfill_imported_senders_regions.php](../../bb_api/backfill_imported_senders_regions.php#L1)
- [bp_earn.php](../../bb_api/bp_earn.php#L1)
- [bp_purchase_request.php](../../bb_api/bp_purchase_request.php#L1)
- [bp_status.php](../../bb_api/bp_status.php#L1)
- [contact_public.php](../../bb_api/contact_public.php#L1)
- [expert_opinion_request.php](../../bb_api/expert_opinion_request.php#L1)
- [expert_request_user_action.php](../../bb_api/expert_request_user_action.php#L1)
- [expert_request_user_cases.php](../../bb_api/expert_request_user_cases.php#L1)
- [expert_request_user_list.php](../../bb_api/expert_request_user_list.php#L1)
- [health.php](../../bb_api/health.php#L1)
- [import_full_datastore_export.php](../../bb_api/import_full_datastore_export.php#L1)
- [index.php](../../bb_api/index.php#L1)
- [login.php](../../bb_api/login.php#L1)
- [me.php](../../bb_api/me.php#L1)
- [messages_list.php](../../bb_api/messages_list.php#L1)
- [messages_mark_read.php](../../bb_api/messages_mark_read.php#L1)
- [migrate_sqlite_to_pgsql.php](../../bb_api/migrate_sqlite_to_pgsql.php#L1)
- [paper_letter_download_link.php](../../bb_api/paper_letter_download_link.php#L1)
- [paper_letter_pdf.php](../../bb_api/paper_letter_pdf.php#L1)
- [paper_letter_send.php](../../bb_api/paper_letter_send.php#L1)
- [partner_click.php](../../bb_api/partner_click.php#L1)
- [partner_portal.php](../../bb_api/partner_portal.php#L1)
- [partners_catalog.php](../../bb_api/partners_catalog.php#L1)
- [pdf.php](../../bb_api/pdf.php#L1)
- [pdf_text.php](../../bb_api/pdf_text.php#L1)
- [question.php](../../bb_api/question.php#L1)
- [register.php](../../bb_api/register.php#L1)
- [reprocess_sender_data.php](../../bb_api/reprocess_sender_data.php#L1)
- [secure_doc_view.php](../../bb_api/secure_doc_view.php#L1)
- [sender_stats.php](../../bb_api/sender_stats.php#L1)
- [social_links.php](../../bb_api/social_links.php#L1)
- [social_links_public.php](../../bb_api/social_links_public.php#L1)
- [solutions_download.php](../../bb_api/solutions_download.php#L1)
- [solutions_download_link.php](../../bb_api/solutions_download_link.php#L1)
- [solutions_list.php](../../bb_api/solutions_list.php#L1)
- [solutions_preview_link.php](../../bb_api/solutions_preview_link.php#L1)
- [stripe_webhook.php](../../bb_api/stripe_webhook.php#L1)
- [support_send.php](../../bb_api/support_send.php#L1)
- [temp_pdf_download.php](../../bb_api/temp_pdf_download.php#L1)
- [template_generate.php](../../bb_api/template_generate.php#L1)
- [tg_bp_bridge.php](../../bb_api/tg_bp_bridge.php#L1)
- [tg_login.php](../../bb_api/tg_login.php#L1)
- [usage.php](../../bb_api/usage.php#L1)
- [user_lang_set.php](../../bb_api/user_lang_set.php#L1)
- [whatsapp_ingest.php](../../bb_api/whatsapp_ingest.php#L1)

### Administration

- [admin/admin_bp_grant.php](../../bb_api/admin/admin_bp_grant.php#L1)
- [admin/admin_message_send.php](../../bb_api/admin/admin_message_send.php#L1)
- [admin/analysis_counter.php](../../bb_api/admin/analysis_counter.php#L1)
- [admin/analysis_requests.php](../../bb_api/admin/analysis_requests.php#L1)
- [admin/analysis_retry.php](../../bb_api/admin/analysis_retry.php#L1)
- [admin/auth_login.php](../../bb_api/admin/auth_login.php#L1)
- [admin/bp_grants_list.php](../../bb_api/admin/bp_grants_list.php#L1)
- [admin/common.php](../../bb_api/admin/common.php#L1)
- [admin/dashboard_stats.php](../../bb_api/admin/dashboard_stats.php#L1)
- [admin/expert_request_assign.php](../../bb_api/admin/expert_request_assign.php#L1)
- [admin/expert_request_delete.php](../../bb_api/admin/expert_request_delete.php#L1)
- [admin/expert_requests_list.php](../../bb_api/admin/expert_requests_list.php#L1)
- [admin/expert_users_stats.php](../../bb_api/admin/expert_users_stats.php#L1)
- [admin/letters.php](../../bb_api/admin/letters.php#L1)
- [admin/paper_letter_deliveries.php](../../bb_api/admin/paper_letter_deliveries.php#L1)
- [admin/parsed_text_purge.php](../../bb_api/admin/parsed_text_purge.php#L1)
- [admin/parsed_text_retention.php](../../bb_api/admin/parsed_text_retention.php#L1)
- [admin/partner_click_stats.php](../../bb_api/admin/partner_click_stats.php#L1)
- [admin/partner_region_request_reply.php](../../bb_api/admin/partner_region_request_reply.php#L1)
- [admin/partner_region_requests.php](../../bb_api/admin/partner_region_requests.php#L1)
- [admin/partners_catalog_add.php](../../bb_api/admin/partners_catalog_add.php#L1)
- [admin/partners_catalog_delete.php](../../bb_api/admin/partners_catalog_delete.php#L1)
- [admin/partners_catalog_list.php](../../bb_api/admin/partners_catalog_list.php#L1)
- [admin/partners_catalog_update.php](../../bb_api/admin/partners_catalog_update.php#L1)
- [admin/partners_list.php](../../bb_api/admin/partners_list.php#L1)
- [admin/partners_photo_upload.php](../../bb_api/admin/partners_photo_upload.php#L1)
- [admin/partners_save.php](../../bb_api/admin/partners_save.php#L1)
- [admin/purchase_mark_paid.php](../../bb_api/admin/purchase_mark_paid.php#L1)
- [admin/purchase_requests.php](../../bb_api/admin/purchase_requests.php#L1)
- [admin/security_alert_resolve.php](../../bb_api/admin/security_alert_resolve.php#L1)
- [admin/security_overview.php](../../bb_api/admin/security_overview.php#L1)
- [admin/senders.php](../../bb_api/admin/senders.php#L1)
- [admin/senders_reprocess.php](../../bb_api/admin/senders_reprocess.php#L1)
- [admin/sent_messages.php](../../bb_api/admin/sent_messages.php#L1)
- [admin/social_links.php](../../bb_api/admin/social_links.php#L1)
- [admin/solutions_delete.php](../../bb_api/admin/solutions_delete.php#L1)
- [admin/solutions_list.php](../../bb_api/admin/solutions_list.php#L1)
- [admin/solutions_preview.php](../../bb_api/admin/solutions_preview.php#L1)
- [admin/solutions_upload.php](../../bb_api/admin/solutions_upload.php#L1)
- [admin/support_messages.php](../../bb_api/admin/support_messages.php#L1)
- [admin/topics.php](../../bb_api/admin/topics.php#L1)
- [admin/users.php](../../bb_api/admin/users.php#L1)
- [admin/users_search.php](../../bb_api/admin/users_search.php#L1)
- [admin/welcome_message.php](../../bb_api/admin/welcome_message.php#L1)

### Buddy Manager

- [buddy-manager/api/auth.php](../../bb_api/buddy-manager/api/auth.php#L1)
- [buddy-manager/api/bewoners.php](../../bb_api/buddy-manager/api/bewoners.php#L1)
- [buddy-manager/api/bootstrap.php](../../bb_api/buddy-manager/api/bootstrap.php#L1)
- [buddy-manager/api/email_send.php](../../bb_api/buddy-manager/api/email_send.php#L1)
- [buddy-manager/api/fax_send.php](../../bb_api/buddy-manager/api/fax_send.php#L1)
- [buddy-manager/api/form_help_analyze.php](../../bb_api/buddy-manager/api/form_help_analyze.php#L1)
- [buddy-manager/api/form_help_fill.php](../../bb_api/buddy-manager/api/form_help_fill.php#L1)
- [buddy-manager/api/form_help_pdf.php](../../bb_api/buddy-manager/api/form_help_pdf.php#L1)
- [buddy-manager/api/generate.php](../../bb_api/buddy-manager/api/generate.php#L1)
- [buddy-manager/api/improve.php](../../bb_api/buddy-manager/api/improve.php#L1)
- [buddy-manager/api/letters.php](../../bb_api/buddy-manager/api/letters.php#L1)
- [buddy-manager/api/pdf.php](../../bb_api/buddy-manager/api/pdf.php#L1)
- [buddy-manager/api/send.php](../../bb_api/buddy-manager/api/send.php#L1)
- [buddy-manager/api/sender_profile.php](../../bb_api/buddy-manager/api/sender_profile.php#L1)
- [buddy-manager/api/upload.php](../../bb_api/buddy-manager/api/upload.php#L1)
- [buddy-manager/api/users.php](../../bb_api/buddy-manager/api/users.php#L1)

### Scripts

- [scripts/build_geo_catalog_from_geonames.php](../../bb_api/scripts/build_geo_catalog_from_geonames.php#L1)

