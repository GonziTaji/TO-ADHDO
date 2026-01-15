import { UnexpectedTagNameError } from "../../utils/errors.js"
import { EVENT_NAMES } from "../../utils/events.js"
import validateSelectors from "../../utils/validate_selectors.js"

const form = {
    init,
}
export default form

const selectors = Object.freeze({
    available_tags_datalist: '#datalist__available_tags',
    suggested_tags_list: '.article_form__suggested_tags_list',
    selected_tags_list: '.article_form__selected_tags_list',
    selected_tag_template: '.article_form__selected_tag_template'
})

function init() {
    /** @type {HTMLFormElement} */
    const form = document.querySelector('#article_form')

    if (form?.tagName !== 'FORM') {
        throw UnexpectedTagNameError(`expected "FORM", got: "${form?.tagName}"`)
    }

    if (form.getAttribute('data-mounted') === 'true') {
        throw AlreadyMountedError(form)
    }

    const error = validateSelectors(form, Object.values(selectors))

    if (error) {
        console.error("Error initializing controller")
        throw error
    }

    form.addEventListener('submit', handleFormSubmit);

    const input = getTagSearchTermInput(form)
    input.addEventListener("input", handleTagInputChange);
    input.addEventListener("keydown", handleTagInputKeyDown);

    form.setAttribute('data-mounted', 'true')
}

function getTagSearchTermInput(form_node) {
    const input = form_node.querySelector('input[name="tag_search_term"]');
    return input
}

/** @type {SubmitEvent} ev */
async function handleFormSubmit(ev) {
    ev.preventDefault();

    /** @type {HTMLFormElement} */
    const form = ev.currentTarget

    if (form.tagName !== 'FORM') {
        throw UnexpectedTagNameError('FORM', form.tagName)
    }

    const form_data = new FormData(form);

    const response = await fetch("/api/articles", { method: "POST", body: form_data })

    const { id, error } = await response.json()

    if (error) {
        alert(error.message || String(error))
        return
    }

    form.reset()

    const suggested_tags_list = form.querySelector(selectors.suggested_tags_list)
    suggested_tags_list.innerHTML = ""

    const selected_tags_list = form.querySelector(selectors.selected_tags_list)
    selected_tags_list.innerHTML = ""

    document.dispatchEvent(new CustomEvent(EVENT_NAMES.new_article, { detail: { id } }))
}

/** @param {KeyboardEvent} ev */
function handleTagInputChange(ev) {
    /** @type {HTMLInputElement} */
    const input = ev.currentTarget
    const form = input.form

    const options_selector = `${selectors.available_tags_datalist} option`
    /** @type {HTMLOptionElement[]} */
    const tags_options = [...form.querySelectorAll(options_selector)]

    const search = input.value.toLowerCase().trim();

    const filtered_tags = !search ? [] : [...tags_options]
        .filter((option) => !option.disabled)
        .filter((option) => option.value.includes(search))
        .map((option) => option.value)
        .sort();

    const ul = document.querySelector(selectors.suggested_tags_list)

    ul.innerHTML = "";

    for (const tag of filtered_tags) {
        const template = document.createElement("template")

        template.innerHTML = `
            <li><button type="button" data-tag-name="${tag}">
              ${tag}
            </button></li>
        `

        const li = template.content.firstChild
        /** @type {HTMLButtonElement}:Outline */
        const btn = li.querySelector('button')

        ul.appendChild(li);

        btn.addEventListener("click", () => {
            addTagToArticle(btn.form, tag)
        })
    }
}

/** @param {KeyboardEvent} ev */
function handleTagInputKeyDown(ev) {
    /** @type {HTMLInputElement} */
    const input = ev.currentTarget;
    const value = input.value.trim().toLowerCase();

    if (value === "") {
        return;
    }

    if (ev.key === "Enter" || ev.key === "Tab") {
        ev.preventDefault();

        addTagToArticle(input.form, value)
    }
}

/**
 * @param {HTMLFormElement} form
 * @param {string} tag_name
 * */
function addTagToArticle(form, tag_name) {
    /** @type {HTMLTemplateElement} */
    const template_node = form.querySelector(selectors.selected_tag_template)
    const new_tag_template_node = template_node.content.cloneNode(true);

    /** @type {HTMLOptionElement} */
    const selected_tag_option = form.querySelector(
        `${selectors.available_tags_datalist} option[value="${tag_name}"]:not(:disabled)`
    )

    const selected_task_list = form.querySelector(selectors.selected_tags_list)
    const tasks_selected = [...selected_task_list.querySelectorAll('input')].map(({ value }) => value)

    if (!tasks_selected.includes(tag_name)) {
        const new_tag_node = new_tag_template_node.children[0]

        const remove_handler = () => {
            new_tag_node.remove();

            if (selected_tag_option) {
                selected_tag_option.disabled = false;
            }
        }

        new_tag_node.querySelector("input").value = tag_name;
        new_tag_node.querySelector("button").addEventListener("click", remove_handler);

        const selected_tags_list = form.querySelector(selectors.selected_tags_list)

        selected_tags_list.appendChild(new_tag_node)
    }

    const suggested_tags_list = form.querySelector(selectors.suggested_tags_list)
    suggested_tags_list.innerHTML = ""

    const input = getTagSearchTermInput(form);
    input.value = ""

    if (selected_tag_option) {
        selected_tag_option.disabled = true;
    }
}

