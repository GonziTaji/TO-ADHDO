document.addEventListener("DOMContentLoaded", init)

function init() {
    bindEvents()
}

function bindEvents() {
    document
        .querySelector('[data-component="delete-dialog"]')
        .addEventListener("close", closeDeleteDialogHandler)

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

    dialog.setAttribute('data-articleid', article_id)
    dialog.showModal()
}

function cancelDeleteHandler() {
    getDeleteDialog().close()
}

function confirmDeleteHandler() {
    getDeleteDialog().close()
}

/** @returns {HTMLDialogElement} */
function getDeleteDialog() {
    return document.querySelector('[data-component="delete-dialog"]')

}

