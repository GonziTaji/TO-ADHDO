
document.addEventListener("DOMContentLoaded", init)

const elements = {
    /** @type {HTMLFormElement */
    form: null,
}

function init() {
    bindEvents()
}

function bindEvents() {
    elements.form = document.querySelector('[data-component="articles-form"]')
    elements.form.addEventListener("submit", handleFormSubmit)
}

/** @param {SubmitEvent} ev */
function handleFormSubmit(ev) {
    ev.preventDefault()

    const data = new FormData(ev.currentTarget)

    const name = data.get("name")
    const description = data.get("description")
    const id = data.get("id")

    console.log(data)
}

