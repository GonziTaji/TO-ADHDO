import { UnexpectedTagNameError } from "../utils/errors.js"
import validateSelectors from "../utils/validate_selectors.js"

const task_form = {
    init,
}
export default task_form

const selectors = Object.freeze({
    available_tags_datalist: '.task_template_form__available_tags_datalist',
    suggested_tags_list: '.task_template_form__suggested_tags_list',
    selected_tags_list: '.task_template_form__selected_tags_list',
    tag_input: 'input[name=tag_selector]',
    selected_tag_template: 'template.task_template_form__selected_tag_template'
})

/** @param {HTMLFormElement} task_form_node */
function init(task_form_node) {
    if (task_form_node.getAttribute('data-mounted') === 'true') {
        throw AlreadyMountedError(task_list_node)
    }

    const error = validateSelectors(task_form_node, Object.values(selectors))

    if (error) {
        console.error("Error initializing task_template_form controller")
        throw error
    }

    task_form_node.addEventListener('submit', handleTaskFormSubmit);

    const input = task_form_node.querySelector(selectors.tag_input);
    input.addEventListener("input", handleTagInputChange);
    input.addEventListener("keydown", handleTagInputKeyDown);

    task_form_node.setAttribute('data-mounted', 'true')
}

/** @type {SubmitEvent} ev */
async function handleTaskFormSubmit(ev) {
    ev.preventDefault();

    /** @type {HTMLFormElement} */
    const form = ev.currentTarget

    if (form.tagName !== 'FORM') {
        throw UnexpectedTagNameError('FORM', form.tagName)
    }

    const form_data = new FormData(form);

    const response = await fetch("/api/tasks_templates", { method: "POST", body: form_data })

    const { error } = await response.json()

    if (error) {
        alert(error.message || String(error))
        return
    }

    form.reset()

    const suggested_tags_list = form.querySelector(selectors.suggested_tags_list)
    suggested_tags_list.innerHTML = ""

    const selected_tags_list = form.querySelector(selectors.selected_tags_list)
    selected_tags_list.innerHTML = ""
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
        .map((option) => option.value);

    const ul = document.querySelector(selectors.suggested_tags_list)

    ul.innerHTML = "";

    for (const tag of filtered_tags) {
        const btn = document.createElement("button");

        btn.type = "button";
        btn.innerText = tag;

        btn.addEventListener("click", () => {
            addTagToTask(tag);
        });

        const li = document.createElement("li");

        li.appendChild(btn);
        ul.appendChild(li);
    }
}

/** @param {KeyboardEvent} ev */
function handleTagInputKeyDown(ev) {
    /** @type {HTMLInputElement} */
    const input = ev.currentTarget;
    const form = input.form
    const value = input.value.trim().toLowerCase();

    if (value === "") {
        return;
    }

    if (ev.key === "Enter" || ev.key === "Tab") {
        ev.preventDefault();

        /** @type {HTMLTemplateElement} */
        const template_node = form.querySelector(selectors.selected_tag_template)
        const new_tag_template_node = template_node.content.cloneNode(true);
        const new_tag_node = new_tag_template_node.children[0]

        /** @type {HTMLOptionElement} */
        const selected_tag_option = form.querySelector(
            `${selectors.available_tags_datalist} option[value="${value}"]`
        )

        const remove_handler = () => {
            new_tag_node.remove();

            if (selected_tag_option) {
                selected_tag_option.disabled = false;
            }
        }

        new_tag_node.querySelector("input").value = value;
        new_tag_node.querySelector("button").addEventListener("click", remove_handler);

        const selected_tags_list = form.querySelector(selectors.selected_tags_list)
        selected_tags_list.appendChild(new_tag_node)

        const suggested_tags_list = form.querySelector(selectors.suggested_tags_list)
        suggested_tags_list.innerHTML = ""

        const input = form.querySelector(selectors.tag_input);
        input.value = ""

        if (selected_tag_option) {
            selected_tag_option.disabled = true;
        }
    }
}

