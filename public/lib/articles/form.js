
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
                removeTag(btn.closest('li'), ds.tagid)
                break
        }
    })

    document
        .querySelector('[data-component="articles-form"]')
        .addEventListener("submit", formSubmitHandler)
}

/**
 * @param {HTMLElement} container
 * @param {string} tag_id
 */
function removeTag(container, tag_id) {
    container.remove()

    if (tag_id) {
        getTagOption(tag_id).disabled = false
    }
}

/** @param {string} tag_id */
function addTag(tag_id) {
    /** @type {HTMLElement} */
    const template = document.querySelector('template[data-component="selected-tag-template"]')
    const new_tag_node = template.content.cloneNode(true).firstChild;

    new_tag_node.querySelector('input[name="tags_names"]').value = getTagOption(tag_id).innerText
    new_tag_node.querySelector('input[name="tags_ids"]').value = tag_id

    document
        .querySelector('[data-component="selected-tags-list"]')
        .appendChild(new_tag_node)

    getTagOption(tag_id).disabled = true
}

/** @param {string} tag_name */
function addNewTag(tag_name) {
    /** @type {HTMLTemplateElement} */
    const template = document.querySelector('template[data-component="selected-tag-template"]')
    const new_tag_node = template.content.cloneNode(true).firstChild;

    new_tag_node.querySelector('input[name="tags_names"]').value = tag_name

    document
        .querySelector('[data-component="selected-tags-list"]')
        .appendChild(new_tag_node)
}

/** @param {SubmitEvent} ev */
function formSubmitHandler(ev) {
    ev.preventDefault()

    const data = new FormData(ev.currentTarget)

    // const name = data.get("name")
    // const description = data.get("description")
    // const id = data.get("id")

    console.log(data)
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

        const suggested_tag_node = template.content.cloneNode(true).firstChild;

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

        const option = getTagOptionByName(value)
        if (option) {
            addTag(option.value)
        } else {
            addNewTag(value)
        }
    }
}

/** 
 * @param {string} tag_id
 * @returns {HTMLOptionElement
 */
function getTagOptionByValue(tag_id) {
    const option = document.querySelector(
        `datalist#datalist-available-tags option[value="${tag_id}"]`
    )

    return option
}

/** 
 * @param {string} tag_name
 * @returns {HTMLOptionElement | undefined}
 */
function getTagOptionByName(tag_name) {
    const option = Array.from(document.querySelectorAll(
        'datalist#datalist-available-tags option'
    ))
        .find(opt => opt.innerText.trim().toLowerCase() == tag_name)

    return option
}
