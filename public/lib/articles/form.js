import { getFirstChildCopyFromTemplate } from "../utils/teststs.js"

document.addEventListener("DOMContentLoaded", init)

function init() {
    bindEvents()
}

function bindEvents() {
    const form = document.querySelector('[data-component="articles-form"]')
    form.addEventListener("submit", formSubmitHandler)

    {
        const tag_search_input = form.querySelector('input[name="tag_search"]')
        tag_search_input.addEventListener("input", tagSearchChangeHandler)
        tag_search_input.addEventListener("keydown", tagSearchKeyDownHandler)
    }

    document.addEventListener('click', (e) => {
        /** @type {HTMLButtonElement} */
        const btn = e.target.closest('button[data-action]')

        if (!btn) return

        const ds = btn.dataset

        switch (ds.action) {
            case 'add-tag':
                addTag(ds.tagid)
                break
            case 'remove-tag':
                removeTag(btn.closest('li'))
                break
        }
    })

    document
        .querySelector('[data-component="articles-form"]')
        .addEventListener("submit", formSubmitHandler)
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
    const form = document.querySelector('[data-component="articles-form"]')
    const fd = new FormData(form)

    if (fd.getAll('tags_names').includes(tag_name)) {
        document.querySelector('#tag_search').value = ""
        return;
    }

    const tag_option = getTagOptionByName(tag_name)

    /** @type {HTMLElement} */
    const template = document.querySelector('template[data-component="selected-tag-template"]')
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

/** @param {string} search_term */
function tagSearchChangeHandler(search_term) {
    const suggested_tags_list = document.querySelector(
        '[data-component="suggested-tags-list"]'
    )

    suggested_tags_list.innerHTML = "";

    /** @type {HTMLOptionElement */
    const tags_options = [...document.querySelectorAll(
        'datalist#datalist-available-tags option'
    )]

    if (!search_term) {
        return
    }

    /** @type {{name: string, id: string}[]} */
    const filtered_tags = tags_options
        .filter((option) => !option.disabled)
        .filter((option) => option.value.includes(search_term))
        .map((option) => ({ name: option.innerText, id: option.value }))

    /** @type {HTMLTemplateElement} */
    const template = document.querySelector(
        'template[data-component="suggested-tag-template"]'
    )

    for (const tag of filtered_tags) {
        const suggested_tag_node = getFirstChildCopyFromTemplate(template)

        /** @type {HTMLButtonElement} */
        const add_button = suggested_tag_node.querySelector('button[data-action="add-tag"]')
        add_button.setAttribute('data-tagid', tag.id)
        add_button.innerText = tag.name

        suggested_tags_list.appendChild(suggested_tag_node);
    }
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
