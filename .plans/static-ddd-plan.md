Plan: Domain-based static assets

Goal
- Move domain-specific static assets into their domain folders and serve them from domain-scoped routes.
- Keep uploads at /public/media/uploads.
- Move shared assets to domain/shared/static and serve from a shared route.
- Use relative paths in templates where possible (static/... under the domain routes).

1) Inventory and mapping
- Enumerate assets under public/ and map each to a target location:
  - Domain assets -> domain/<domain>/static/...
  - Shared assets -> domain/shared/static/...
  - Uploads remain in public/media/uploads
- Find all template references to /public/... and any hardcoded paths in JS/CSS.

2) Define routing strategy
- Keep a narrow global static handler for:
  - /favicon.ico
  - /public/media/uploads (uploads only)
- Add per-domain static mounts inside each domain's RegisterRoutes:
  - Example (articles):
    - group.Static("/static", "domain/articles/static")
    - admin.Static("/static", "domain/articles/static")
- Add a shared static mount in server startup:
  - router.Static("/shared/static", "domain/shared/static")

3) Move files
- Domain-specific assets:
  - public/lib/articles/form.js -> domain/articles/static/js/form.js
  - public/lib/articles/list.js -> domain/articles/static/js/list.js
  - public/media/spinner.svg -> domain/articles/static/media/spinner.svg
- Shared assets:
  - public/lib/utils/* -> domain/shared/static/js/*
  - public/styles/main.css -> domain/shared/static/css/main.css
- Leave public/media/uploads/** untouched.

4) Update template paths (relative where possible)
- Articles templates:
  - src="static/js/form.js"
  - src="static/js/list.js"
  - src="static/media/spinner.svg"
- Shared CSS (if used):
  - href="/shared/static/css/main.css"

5) Update JS/CSS internal references
- Search moved JS/CSS for any /public/... references and update to new routes.

6) Fix/verify article list template usage
- domain/articles/views.go references public/lib/components/articles_list/template.html which does not exist.
  - Decide whether to add/move this template under domain/articles/static/templates or remove the dependency if unused.
  - Update ParseFiles path accordingly.

7) Update server/static.go
- Reduce registerStaticRoutes to only global assets:
  - /favicon.ico
  - /public/media/uploads
  - /shared/static
- Domain-specific static mounts live in each domain's RegisterRoutes.

8) Smoke check
- Start server and verify:
  - /articles and /admin/articles/* pages load JS/CSS/media without 404s.
  - Uploads still resolve at /public/media/uploads/...
  - CSP allows the new paths (should be covered by 'self').
