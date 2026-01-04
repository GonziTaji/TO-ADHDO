import elements from "../utils/elements.js"

document.addEventListener("DOMContentLoaded", init);

function init() {
    elements.task_form().addEventListener('submit', handleTaskFormSubmit);

    const delete_task_btns = elements.task_list_container().querySelectorAll('.delete_task_btn');
    delete_task_btns.forEach((btn) => btn.addEventListener('click', handleDeleteTaskClick));
}

/** @type {SubmitEvent} ev */
async function handleTaskFormSubmit(ev) {
    ev.preventDefault();

    const form_data = new FormData(ev.currentTarget);

    const response = await fetch("/api/tasks", { method: "POST", body: form_data })

    const { new_tags } = await response.json()

    for (const new_tag of new_tags) {
        const option = document.createElement('option');
        option.value = new_tag;

        elements.tags_options_list().appendChild(option);
    }

    elements.filtered_tags_list().innerHTML = "";
    elements.selected_tags_list().innerHTML = "";
    elements.task_name_input().value = "";
    elements.task_description_textarea().value = "";

    const task_list_response = await fetch("/api/tasks")

    elements.task_list_container().innerHTML = await task_list_response.text();
}

/** @param {MouseEvent} ev */
async function handleDeleteTaskClick(ev) {
    /** @type {HTMLButtonElement} */
    const btn = ev.currentTarget;
    const task_id = btn.getAttribute('data-task-id');

    fetch(`/api/task${task_id}`, { method: "DELETE" })

    elements.task_list_container().querySelector(`.task_list_item[data-task-id="${task_id}"]`).remove()
}

