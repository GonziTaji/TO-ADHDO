document.addEventListener("DOMContentLoaded", init)

function init() {
    bindEvents()
}

function bindEvents() {
    document.addEventListener('beforetoggle', (/** @type {ToggleEvent} */ ev) => {
        /** @type {HTMLDialogElement | null} */
        const dialog = ev.target.closest('dialog')

        if (!dialog) {
            return;
        }

        dialog.dataset.state = newState
    })

    const delete_dialog = getDeleteDialog()
    delete_dialog.addEventListener("close", closeDeleteDialogHandler)

    document.addEventListener('click', (ev) => {
        /** @type {HTMLButtonElement | null} */
        const btn = ev.target.closest('button[data-action]')

        if (!btn) return

        switch (btn.dataset.action) {
            case "new":
                newArticleHandler()
                break

            case "view":
                viewArticleHandler(btn.dataset.articleid)
                break

            case "edit":
                editArticleHandler(btn.dataset.articleid)
                break

            case "delete":
                deleteArticleHandler(btn.dataset.articleid)
                break

            case "cancel-delete":
                cancelDeleteHandler()
                break

            case "confirm-delete":
                confirmDeleteHandler()
                break

            default:
                console.warn("Invalid button action: " + btn.dataset.action)
        }
    })

}

/** @param {Event} ev */
function closeDeleteDialogHandler(ev) {
    delete ev.currentTarget.dataset.articleid
}

/** @param {KeyboardEvent} _ */
function newArticleHandler(_) {
    location.href = "new"
}

/** param {string} article_id */
function viewArticleHandler(article_id) {
    location.href = `${article_id}`
}

/** param {string} article_id */
function editArticleHandler(article_id) {
    location.href = `${article_id}/edit`
}

/** param {string} article_id */
function deleteArticleHandler(article_id) {
    const dialog = getDeleteDialog()

    dialog.dataset.articleid = article_id
    dialog.showModal()
}

function cancelDeleteHandler() {
    const dialog = getDeleteDialog()

    delete dialog.dataset.articleid

    getDeleteDialog().close()
}

async function confirmDeleteHandler() {
    const dialog = getDeleteDialog()
    const article_id = dialog.dataset.articleid;

    const item = document.querySelector(
        `[data-component="articles-list-item"]:has([data-articleid="${article_id}"])`
    )

    item.dataset.loading = true

    const response = await fetch(`/admin/articles/${article_id}`, { method: 'DELETE' })

    item.dataset.loading = false
    dialog.close()

    if (!response.ok) {
        alert("Uh oh!" + await response.text())
        return
    }

    // Should we ask the backend for the full list to ensure the list is updated with the data?
    // It would guard against false positive creations responses
    item.remove()

    alert("Article deleted successfully")
}

/** @returns {HTMLDialogElement} */
function getDeleteDialog() {
    return document.querySelector('#article-list-delete-article-dialog')

}

