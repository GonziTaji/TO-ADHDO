import postWithProgress from '../utils/postWithProgress.js'
import { getFirstChildCopyFromTemplate } from "../utils/teststs.js"

document.addEventListener("DOMContentLoaded", init)

const MAX_ARTICLE_IMAGES = 5

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
            case 'image_uploader':
                if (!input.files && input.files.length === 0) {
                    return
                }

                const existing_image_row = input.closest('[data-imageid]')

                if (existing_image_row) {
                    existing_image_row.remove()
                }

                const files = [...input.files]
                input.value = ""

                files.forEach(uploadImage)

                break
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

            case 'remove-image':
                removeImage(btn.closest('[data-imageid]'))
                break
        }
    })
}

/** @param {HTMLElement} container */
function removeImage(container) {
    if (!container) {
        console.error(new Error('no container to remove'))
        return
    }

    container.remove()
}

function bindDragAndDropEvents() {
    document.addEventListener('dragover', (e) => {
        e.preventDefault()
    })

    document.addEventListener('dragleave', (e) => {
        /** @type {HTMLElement} */
        const subject = e.target;

        if (subject.dataset.is_drag_over !== undefined && !subject.contains(e.fromElement)) {
            subject.dataset.is_drag_over = ''
        }
    })

    document.addEventListener('dragenter', (e) => {
        /** @type {HTMLElement} */
        const subject = e.target.closest('[data-is_drag_over]');

        if (subject) {
            subject.dataset.is_drag_over = 'true'
        }
    })

    document.addEventListener("drop", (e) => {
        e.preventDefault()

        document.querySelectorAll('[data-is_drag_over]').forEach((node) => {
            node.dataset.is_drag_over = ''
        })

        /** @type {HTMLElement} */
        const subject = e.target.closest('[data-is_drag_over]')

        if (!subject) {
            return
        }

        if (subject.dataset.dropzone_id === "article_images") {
            const fd = new FormData(subject.form)

            const current_images_count = fd.getAll('articles_images').length

            const max_files_to_drop = MAX_ARTICLE_IMAGES - current_images_count

            console.log(`max files to drop: ${max_files_to_drop}`)

            if (max_files_to_drop === 0) {
                console.log('max files reached')
                return
            }

            [...e.dataTransfer.files]
                .filter((file) => file.type.startsWith("image/"))
                .splice(0, max_files_to_drop)
                .forEach(uploadImage)
        }
    })
}

/**
 * @param {string} title
 * @param {string} content
 */
function showErrorModal(title, content) {
    const fallback = (/** @type {string} */ el_name) => {
        console.error(new Error(`Error: "${el_name}" element could not be found`))
        alert(`${title}. ${content}`)
    }

    /** @type {HTMLDialogElement | null} */
    const dialog = document.querySelector('dialog#error-dialog')

    if (!dialog) {
        return fallback('dialog')
    }

    const title_el = dialog.querySelector('#error-dialog_title')
    const content_el = dialog.querySelector('#error-dialog_content')

    if (!title_el || !content_el) {
        return fallback("title or content")
    }

    title_el.innerText = title
    content_el.innerText = content

    dialog.showModal()
}

/** @param {File} file */
async function uploadImage(file) {
    const template = document.getElementById('articles/form/image-miniature--loader')
    const image_loader = getFirstChildCopyFromTemplate(template)

    const url = URL.createObjectURL(file)

    image_loader.querySelector('img[src=""]').src = url

    const images_grid = document.getElementById('images-grid')
    images_grid.appendChild(image_loader)

    const fd = new FormData()
    fd.set('file', file)

    const error_title = 'Error uploading image'

    try {
        const response = await postWithProgress('/articles/uploads', fd, (progress => {
            image_loader.dataset.progress = progress
        }))

        if (!response.ok) {
            console.log(response)
            showErrorModal(error_title, `${response.statusText}: ${response.body}`)
            return
        }

        URL.revokeObjectURL(url)

        const dest = document.createElement('div')

        images_grid.appendChild(dest)

        dest.outerHTML = response.body
    } catch (e) {
        console.log(e)
        showErrorModal(error_title, `Unexpected error: ${e.message || e}`)
    } finally {
        image_loader.remove()
    }

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
    const article_id = body.get("id")

    let method = 'POST'
    let endpoint = "/articles"

    if (article_id != "") {
        method = 'PUT'
        endpoint += "/" + article_id
    }

    const response = await fetch(endpoint, { method, body })

    if (!response.ok) {
        showErrorModal('Uh oh! ewrorw', await response.text())
        return
    }

    /** @type {HTMLDialogElement | null} */
    const dialog = document.querySelector('dialog#success-dialog')
    /** @type {HTMLAnchorElement | null} */
    const anchor = dialog?.querySelector('a#success-dialog_article_link')

    if (!dialog || !anchor) {
        console.warn('dialog or anchor of the dialog not found in the DOM. using fallback')
        alert('Success!')
        location.href = response.headers.get('location')
        return
    }

    anchor.href = response.headers.get('location')

    dialog.showModal()
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
 * @returns {HTMLOptionElement}
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
