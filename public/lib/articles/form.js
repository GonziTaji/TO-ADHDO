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

            case 'confirm-price':
                confirmPrice()
                break

            case 'reset-price':
                resetPrice()
                break
        }
    })

    document.addEventListener('input', (e) => {
        const input = e.target.closest('input')

        switch (input.id) {
            case 'new_price':
                updateAndReportValidity(input, priceValidator)
                break;
        }
    })
}

/**
 * @param {HTMLInputElement} input
 * @param {(string) => string | null } validator
 */
function updateAndReportValidity(input, validator) {
    const err = validator(input.value)

    input.setCustomValidity(err || '')

    if (err) {
        input.reportValidity()
    }
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

function confirmPrice() {
    const input_selector = 'input#new_price'
    const /** @type {HTMLInputElement} */ input = document.querySelector(input_selector)

    if (!input) {
        console.error(new Error(`no input element matched selector : ${input_selector}`))
        alert('uh oh, fatal error')
        return
    }

    updateAndReportValidity(input, priceValidator)

    if (!input.validity.valid) {
        return
    }

    const loader = document.querySelector('.prices-grid .loader')
    loader.dataset.show = true
}

function resetPrice() {
    const input_selector = 'input#new_price'
    const /** @type {HTMLInputElement} */ input = document.querySelector(input_selector)

    input.disabled = false
    input.value = ''
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
