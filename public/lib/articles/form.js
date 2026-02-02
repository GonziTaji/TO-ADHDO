import { getFirstChildCopyFromTemplate } from "../utils/teststs.js"

document.addEventListener("DOMContentLoaded", init)

function init() {
    bindEvents()
}

function bindEvents() {
    {
        const form = document.getElementById('articles-form')
        form.addEventListener("submit", formSubmitHandler)
        form.addEventListener("reset", formResetHandler)

        const tag_search_input = form.querySelector('input[name="tag_search"]')
        tag_search_input.addEventListener("keydown", tagSearchKeyDownHandler)
    }

    bindDragAndDropEvents()

    document.addEventListener('change', async (e) => {
        /** @type {HTMLInputElement} */
        const input = e.target.closest('input')

        if (!input) return

        switch (input.name) {
            case 'article_image':
                const fd = new FormData()

                if (!input.files || input.files.length == 0) {
                    console.error(new Error('no file'))
                    return
                }

                fd.set('article_image', input.files[0])

                const r = await fetch('/api/uploads/articles_images', { method: 'POST' })

                console.log(await r.text())

                break;
        }
    })

    document.addEventListener('close', (e) => {
        /** @type {HTMLDialogElement} */
        const dialog = e.target.closest('dialog')

        if (!dialog) return

        switch (dialog.id) {
            case 'new-price-dialog':
                handleNewPriceSubmit(dialog)
                break
        }
    }, { capture: true })

    document.addEventListener('click', (e) => {
        /** @type {HTMLButtonElement} */
        const btn = e.target.closest('button[data-action]')

        if (!btn) return

        const { action, tagid } = btn.dataset

        switch (action) {
            case 'add-tag':
                addTag(tagid)
                break
            case 'remove-tag':
                removeTag(btn.closest('li'))
                break

            case 'remove-new-price':
                removeNewPrice()
                break
        }
    })
}

function bindDragAndDropEvents() {
    document.addEventListener('dragenter', (e) => {
        const subject = e.target.closest('[data-dragstatus]')
        if (!subject) return
        subject.dataset.dragstatus = e.type
    })

    document.addEventListener("drop", (e) => {
        e.preventDefault()

        const subject = e.target.closest('[data-dragstatus]')

        if (!subject) return
        subject.dataset.dragstatus = ""

        if (subject.id === "article_image_dropzone") {
            /** @type {HTMLInputElement} */
            const file_input = document.querySelector('input[type="file"]#article_image')

            if (!file_input) {
                console.log(new Error('No input element found'))
                return
            }

            file_input.files.ok
        }

    }, { capture: true })

    document.addEventListener('dragover', (e) => {
        e.preventDefault()
    })

    document.addEventListener('dragleave', (e) => {
        const subject = e.target.closest('[data-dragstatus]')
        if (!subject) return
        subject.dataset.dragstatus = e.type
    })

}

/** @param {HTMLDialogElement} dialog */
function handleNewPriceSubmit(dialog) {
    if (!dialog.returnValue) return

    const prices_grid = document.getElementById('prices-grid')

    if (!prices_grid) {
        console.error(new Error('price grid element not found'))
        return
    }

    const grid_price_input = prices_grid.querySelector('[name="new_price"]')
    const grid_description_input = prices_grid.querySelector('[name="new_price_description"]')

    if (!grid_price_input || !grid_description_input) {
        console.error(new Error('one or more elements could not be found'), {
            elements: { grid_price_input, grid_description_input }
        })
        return
    }

    const dialog_form = dialog.querySelector('form')

    if (!dialog_form) {
        console.error(new Error('dialog doesn\'t have a form element'))
        return
    }

    const fd = new FormData(dialog_form)
    const raw_price = fd.get('price')
    const price_err = priceValidator(raw_price)

    if (price_err) {
        const dialog_price_input = dialog_form.querySelector('[name="price"')
        dialog_price_input.setCustomValidity(price_err)
        dialog_price_input.reportValidity()
        return
    }

    // To remove leading zeroes
    const new_price_value = Number(raw_price).toString()

    grid_price_input.value = new_price_value
    grid_description_input.value = fd.get('description').trim()

    prices_grid.dataset.has_new_price = true
}

/**
 * @param {string} price
 * @returns {string | null}
 * */
function priceValidator(price) {
    if (!price || String(price).trim().length == 0) {
        return 'Please set a value for the new price'
    }

    const parsed_price = Number(price)

    if (isNaN(parsed_price)) {
        return 'Please input a valid number'
    }

    if (parsed_price < 1) {
        return 'The price must be greater than 0'
    }

    return null
}

function removeNewPrice() {
    const el = document.querySelector('[data-has_new_price="true"]')

    if (!el) {
        console.error('No match for element to remove new price')
        return
    }

    el.dataset.has_new_price = false

}

/**
 * @param {HTMLElement} tag_container
 * @param {string} tag_id
 */
function removeTag(tag_container) {
    const tag_id = tag_container.querySelector('input[name="tags_ids"]').value

    tag_container.remove()

    if (tag_id) {
        getTagOptionById(tag_id).disabled = false
    }
}

/** @param {string} tag_name */
function addTag(tag_name) {
    const form = document.getElementById('articles-form')
    const fd = new FormData(form)

    if (fd.getAll('tags_names').includes(tag_name)) {
        document.querySelector('#tag_search').value = ""
        return;
    }

    const tag_option = getTagOptionByName(tag_name)

    /** @type {HTMLElement} */
    const template = document.getElementById('selected-tag-template')
    const new_tag_node = getFirstChildCopyFromTemplate(template)

    new_tag_node.querySelector('input[name="tags_names"]').value = tag_name

    if (tag_option) {
        tag_option.disabled = true
        new_tag_node.querySelector('input[name="tags_ids"]').value = tag_option.dataset.tagid
    } else {
        new_tag_node.querySelector('input[name="tags_ids"]').value = ""
    }

    document
        .querySelector('[data-component="selected-tags-list"]')
        .appendChild(new_tag_node)

    document.querySelector('#tag_search').value = ""
}

function formResetHandler() {

}

/** @param {SubmitEvent} ev */
async function formSubmitHandler(ev) {
    ev.preventDefault()

    const body = new FormData(ev.currentTarget)

    let endpoint = "/articles"

    const article_id = body.get("id")

    let method = "POST"

    if (article_id != "") {
        endpoint += "/" + article_id
        method = "PUT"
    }

    const response = await fetch(endpoint, { method, body })

    if (!response.ok) {
        alert('Uh oh! ewrorw:' + await response.text())
        return
    }

    console.log(await response.text())
}

/** @param {KeyboardEvent} ev */
function tagSearchKeyDownHandler(ev) {
    /** @type {HTMLInputElement} */
    const input = ev.currentTarget;
    const value = input.value.trim().toLowerCase();

    if (value === "") {
        return;
    }

    if (ev.key === "Enter" || ev.key === "Tab") {
        ev.preventDefault();

        addTag(value)
    }
}

/** 
 * @param {string} tag_id
 * @returns {HTMLOptionElement
 */
function getTagOptionById(tag_id) {
    const option = document.querySelector(
        `datalist#datalist-available-tags option[data-tagid="${tag_id}"]`
    )

    return option
}

/** 
 * @param {string} tag_name
 * @returns {HTMLOptionElement | undefined}
 */
function getTagOptionByName(tag_name) {
    const option = document.querySelector(
        `datalist#datalist-available-tags option[value="${tag_name}"]`
    )

    return option
}
