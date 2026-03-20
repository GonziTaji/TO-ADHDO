Plan: Connect webapp pages to server controllers

Goal
- Split current mixed HTML/controller routes into clear page routes and API routes.
- Serve the React multipage app from server page routes like `/catalog`, `/catalog/:id`, `/wishlist`, etc.
- Expose JSON API endpoints under `/api/...` for the webapp to consume.
- Keep all existing stores unchanged.
- Do not add cart routes.

Constraints
- Only controllers, routes, server static/page serving, and web API clients/pages should change.
- Stores in `domain/articles/store.go`, `domain/wishlist/store.go`, and `domain/tags/store.go` must remain untouched.
- Cart stays frontend-only; no backend routes or controller work for cart.

1) Final route contract
- Introduce a clear separation between page routes and API routes.
- Canonical page routes:
  - `GET /` -> web home page
  - `GET /catalog` -> public catalog page
  - `GET /catalog/:article_id` -> public article detail page
  - `GET /articles` -> article admin list page
  - `GET /articles/new` -> article create page
  - `GET /articles/:article_id/edit` -> article edit page
  - `GET /wishlist` -> wishlist page
  - `GET /tags` -> tags page
- Canonical API routes:
  - `GET /api/catalog` -> public catalog list/filter data
  - `GET /api/articles` -> admin article list data
  - `GET /api/articles/:article_id` -> article detail data
  - `POST /api/articles` -> create article
  - `PUT /api/articles/:article_id` -> update article
  - `DELETE /api/articles/:article_id` -> delete article
  - `POST /api/articles/uploads` -> upload article image
  - `GET /api/wishlist` -> wishlist data plus filter metadata
  - `GET /api/wishlist/preview` -> existing preview endpoint
  - `GET /api/tags` -> tags list data
- Explicitly out of scope for this migration:
  - any backend cart route
  - new wishlist admin page routes not already represented in the webapp
  - new tag mutation UI routes in the webapp
- Keep temporary compatibility redirects from old paths:
  - `/catalog/` -> `/catalog`
  - `/admin/articles/` -> `/articles`
  - `/admin/articles/new` -> `/articles/new`
  - `/admin/articles/:article_id/edit` -> `/articles/:article_id/edit`

2) Separate page controllers from API controllers
- Articles currently mix HTML page rendering and mutations in `domain/articles/controller_public.go` and `domain/articles/controller_admin.go`.
- Refactor controller responsibilities into two groups without changing service/store behavior:
  - Page handlers: return the built webapp HTML files.
  - API handlers: return JSON for list/detail/create/update/delete/upload.
- For articles:
  - Replace current HTML page handlers (`GetCatalogHandler`, `GetHandler`, `GetListHandler`, `GetFormHandler`) with page-serving handlers or new dedicated page handlers.
  - Reuse existing service calls inside new JSON handlers.
- For wishlist:
  - Split `GetListHandler` into:
    - page handler for `/wishlist`
    - API handler for `/api/wishlist`
  - Keep `GetPreview` as API-only.
- For tags:
  - Split `GetListHandler` into:
    - page handler for `/tags`
    - API handler for `/api/tags`

3) Decide how the server serves the built webapp pages
- Serve Vite build output from `web/dist` as static page assets.
- Add a small server-side page delivery layer that maps page routes to specific built HTML files.
- Recommended mapping:
  - `/` -> `web/dist/index.html`
  - `/catalog` -> `web/dist/articles/catalog.html`
  - `/catalog/:article_id` -> also serve `web/dist/articles/view.html`
  - `/articles` -> `web/dist/articles/list.html`
  - `/articles/new` -> `web/dist/articles/form.html`
  - `/articles/:article_id/edit` -> `web/dist/articles/form.html`
  - `/wishlist` -> `web/dist/wishlist/wishlist.html`
  - `/tags` -> `web/dist/tags/list.html`
- Also mount the Vite assets directory so hashed JS/CSS files load correctly:
  - `/assets/*` -> `web/dist/assets/*`
- Prefer explicit page handlers over a catch-all SPA fallback, since this is a multipage app.

4) Add a reusable server helper for built page responses
- Create a small helper in `server/` for serving HTML files from `web/dist`.
- The helper should:
  - resolve the file path once from repo root / working dir assumptions
  - return 404 if the target file is missing
  - set `Content-Type: text/html; charset=utf-8`
- This avoids duplicating `ctx.File(...)` logic in every domain.

5) Update domain route registration
- `domain/articles/routes.go`
  - Remove page ownership from the old `/admin/articles` and `/catalog/*` template routes.
  - Register new page routes:
    - `/catalog`
    - `/catalog/:article_id`
    - `/articles`
    - `/articles/new`
    - `/articles/:article_id/edit`
  - Register new API routes under `/api/articles` and `/api/catalog`.
- `domain/wishlist/routes.go`
  - Register page route `/wishlist`.
  - Register API routes under `/api/wishlist`.
  - Remove admin page routes from the implementation plan unless they must remain for compatibility; if kept, make them redirects only.
- `domain/tags/routes.go`
  - Register page route `/tags`.
  - Register API routes under `/api/tags`.

6) Define the JSON response contracts before coding
- The current web client expects JSON shapes that do not yet exist on the server.
- Before implementation, define DTOs per endpoint so controller responses are stable and explicit.
- Articles API DTOs should cover:
  - list item: `Id`, `Name`, `Price`, `ThumbnailUrl`, `Tags`, `Prices`, `AvailableForTrade`, optional `Condition`
  - detail item: existing detail fields used by `view.tsx`
  - form payload for create/update compatible with current service input
- Wishlist API DTO should combine the data currently rendered in template `wishlist.html`:
  - `Items`
  - `SearchTerm`
  - `PriceRange`
  - `PriceSelectedRange`
  - `TagsSelectOptions`
- Tags API DTO should return the list currently used by the page.
- Do not expose raw DB/store internals if a smaller response shape is enough.

7) Adapt article controllers to JSON, without touching stores
- Reuse current service methods:
  - `Catalog(...)`
  - `CatalogList(...)`
  - `GetDetails(...)`
  - `List(...)`
  - `GetFormData(...)`
  - `CreateFromForm(...)`
  - `UpdateFromForm(...)`
  - `Delete(...)`
- Add JSON handlers that translate service output to the web DTOs.
- Important normalization work:
  - create a single article detail JSON mapper
  - create a single article list JSON mapper
  - decide whether admin article list should come from `List(...)` and catalog from `CatalogList(...)`, or whether one shared mapper can support both
- Preserve upload behavior by changing only the response format if needed.
  - If the web form still needs image upload, decide whether upload should return JSON instead of the current HTML fragment.
  - Since the current web form does not yet implement uploads, upload migration can be planned as a later step, but the API route should still be defined.

8) Adapt wishlist controllers to JSON, without touching stores
- Reuse `GetWishlist(...)`, `GetWishitem(...)`, `SaveWishitem(...)`, and `DeleteWishitem(...)` through the existing controller/store flow.
- Replace the frontend assumption of `/api/wishlist/filters` with a single `GET /api/wishlist` response that includes both filter metadata and items.
- Keep `GET /api/wishlist/preview` unchanged in spirit, only relocated/confirmed under the new route layout.
- If wishlist item detail/edit pages are not part of the current webapp, do not add new page routes for them in this phase.

9) Adapt tags controllers to JSON, without touching stores
- Reuse `store.List(...)` and `store.Delete(...)` via controller methods.
- Add:
  - `GET /tags` page route -> serves web page
  - `GET /api/tags` -> JSON list for the React page
  - optionally `DELETE /api/tags/:tagid` if the webapp is expected to support delete later
- If delete is not used by the webapp yet, it can remain out of scope for the first pass.

10) Connect page routes to the built web pages
- Update the webapp to assume server routes without `.html` suffixes.
- Replace hardcoded links in `web/src/shared/components/Layout.tsx`:
  - `/articles/catalog.html` -> `/catalog`
  - `/articles/list.html` -> `/articles`
  - `/wishlist/wishlist.html` -> `/wishlist`
  - `/tags/list.html` -> `/tags`
- Update page navigation in article pages:
  - list -> new: `/articles/new`
  - list -> edit: `/articles/:id/edit`
  - list -> detail: `/catalog/:id`
  - detail -> wishlist: `/wishlist`
- Keep cart links unchanged on the frontend, but do not add server planning for cart routes.

11) Update frontend API client paths
- Revise `web/src/api/client.ts` to match the new API surface.
- Recommended frontend API mapping:
  - `articlesApi.list()` -> `GET /api/articles`
  - `articlesApi.get(id)` -> `GET /api/articles/:id`
  - `articlesApi.create(data)` -> `POST /api/articles`
  - `articlesApi.update(id, data)` -> `PUT /api/articles/:id`
  - `articlesApi.delete(id)` -> `DELETE /api/articles/:id`
  - `catalogApi.list(filters)` or `articlesApi.catalog(filters)` -> `GET /api/catalog`
  - `wishlistApi.list(filters)` -> `GET /api/wishlist`
  - remove `wishlistApi.filters()` since wishlist filter metadata should come from the same endpoint
  - `tagsApi.list()` -> `GET /api/tags`
- Update pages accordingly:
  - `catalog.tsx` should read from `/api/catalog`, not `/api/articles`
  - `Wishlist.tsx` should stop requesting `/api/wishlist/filters`

12) Pass route params cleanly to the React pages
- Because `view.html` and `form.html` are reused for dynamic server routes, the React pages should read identifiers from the URL path/query the same way the server exposes them.
- Recommended behavior:
  - `/catalog/:article_id` -> `view.tsx` reads the last path segment as article id
  - `/articles/:article_id/edit` -> `form.tsx` reads the path segment instead of relying on `?id=`
  - `/articles/new` -> `form.tsx` treats missing path id as create mode
- This lets server page routes be canonical and avoids `.html?id=...` style URLs.

13) Remove old template rendering dependencies incrementally
- Once page routes serve built React HTML, the old template-based page handlers become obsolete for those pages.
- Keep template rendering only where still needed by non-web flows, if any.
- Safe migration order:
  1. Add new page routes and JSON API routes.
  2. Update webapp links/API calls.
  3. Verify React pages work end-to-end.
  4. Only then remove old page template routes or convert them to redirects.

14) Backward compatibility and redirects
- To reduce breakage while migrating bookmarks/internal links, consider temporary redirects:
  - `/catalog/` -> `/catalog`
  - `/catalog/:id` can remain if already close to target
  - `/admin/articles/` -> `/articles`
  - `/admin/articles/new` -> `/articles/new`
  - `/admin/articles/:id/edit` -> `/articles/:id/edit`
  - `/tags/` -> `/tags`
  - `/wishlist/` -> `/wishlist`
- Keep these as short-lived migration aids, not permanent architecture.

15) Verification checklist for the later implementation
- Server-side:
  - page routes return HTML and load `/assets/*` correctly
  - API routes return JSON with expected status codes
  - no handler still depends on HTML templates for migrated pages
- Frontend:
  - home page links use clean server routes
  - catalog loads via `/api/catalog`
  - article detail loads via `/api/articles/:id`
  - article create/update/delete use `/api/articles...`
  - wishlist loads from one `/api/wishlist` payload
  - tags load from `/api/tags`
- Data integrity:
  - no store file changed
  - create/update/delete behavior still uses existing service/store paths

16) Suggested implementation order
- Step 1: add server helper for serving `web/dist` pages and static `/assets`
- Step 2: add new page routes for catalog/articles/wishlist/tags
- Step 3: add JSON API handlers and routes for articles, catalog, wishlist, tags
- Step 4: update `web/src/api/client.ts`
- Step 5: update web page links/navigation to clean routes
- Step 6: update dynamic page param parsing in `view.tsx` and `form.tsx`
- Step 7: test all page/API flows against the running Go server
- Step 8: remove or redirect obsolete HTML/template routes

Resolved decisions
- Public article detail lives only at `/catalog/:article_id`.
- `/articles/*` is reserved for admin CRUD pages backed by the webapp.
- Wishlist is exposed only as `/wishlist` page plus `/api/wishlist` and `/api/wishlist/preview` in this phase.
- Tags are exposed only as `/tags` page plus `/api/tags` in this phase.
- Cart remains frontend-only and is excluded from server routing work.
- Article upload endpoint should be placed at `/api/articles/uploads`; response format can be normalized to JSON during implementation if the current React form needs it.
