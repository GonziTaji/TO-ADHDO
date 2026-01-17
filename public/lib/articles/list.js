document.addEventListener("DOMContentLoaded", init)

const elements = {
    /** @type {HTMLDivElement */
    list: null,
    /** @type {HTMLDialogElement */
    dialog: null,
}

function init() {
    bindEvents()
}

function bindEvents() {
    elements.list = document.querySelector('[data-component="articles-list"]')

    elements.list
        .querySelector('button[data-action="new"]')
        .addEventListener("click", handleNewArticleBtnClick)

    elements.list
        .querySelectorAll('[data-component="articles-list-item"] button[data-article-id]')
        .forEach(btn => btn.addEventListener("click", handleArticleActionBtnClick))

    elements.dialog = document.querySelector('[data-component="confirm-delete-article-dialog"]')

    elements.dialog
        .addEventListener("close", ev => ev.currentTarget.removeAttribute('data-article-id'))

    elements.dialog
        .querySelector('[data-action="cancel"]')
        .addEventListener("click", handleCancelDeleteBtnClick)

    elements.dialog
        .querySelector('[data-action="confirm"]')
        .addEventListener("click", handleConfirmDeleteBtnClick)

}

/** @param {KeyboardEvent} ev */
function handleNewArticleBtnClick() {
    location.href = "new"
}

function handleCancelDeleteBtnClick() {
    elements.dialog.close()
}

function handleConfirmDeleteBtnClick() {
    elements.dialog.close()
}


/** @param {KeyboardEvent} ev */
function handleArticleActionBtnClick(ev) {
    /** @type {HTMLButtonElement} */
    const btn = ev.currentTarget
    const action = btn.getAttribute("data-action")
    const article_id = btn.getAttribute("data-article-id")

    if (!action || !article_id) {
        console.warn("Handler called from a button without all the required props. Expected: data-action, data-article-id")
        return
    }

    switch (action) {
        case "view":
            viewArticleHandler(article_id)
            break;

        case "edit":
            editArticleHandler(article_id)
            break;

        case "delete":
            deleteArticleHandler(article_id)
            break;

        default:
            console.error("Unexpected button action: " + action)
    }
}

/** param {string} article_id */
function viewArticleHandler(article_id) {
    location.href = `${article_id}/details`
}

/** param {string} article_id */
function editArticleHandler(article_id) {
    location.href = `${article_id}/form`
}

/** param {string} article_id */
function deleteArticleHandler(article_id) {
    elements.dialog.setAttribute('data-article-id', article_id)
    elements.dialog.showModal()
}
