Small plan (helper-based approach)

- Add a template helper (e.g., `staticBase`) that returns `/admin/articles/static` or `/articles/static` based on the request path.
- Inject the request path into template data from handlers (recommended: a base template data struct that includes `RequestPath`), so the helper can choose the correct base.
- Update templates to use `{{ staticBase .RequestPath }}/js/form.js`, `{{ staticBase .RequestPath }}/js/list.js`, `{{ staticBase .RequestPath }}/media/spinner.svg`.
- Verify both nested admin routes (`/admin/articles/1/edit`) and public routes (`/articles/1`) return 200 for assets.
