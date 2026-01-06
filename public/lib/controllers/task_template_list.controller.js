import { AlreadyMountedError } from "../utils/errors.js";
import service from "./task_template_list.service.js";

const task_templates_list = {
    init,
}

export default task_templates_list

const selectors = {
    delete_task_btn: '.task_template_list_item__delete_btn',
    view_task_btn: '.task_template_list_item__view_btn',
}

/** @param {HTMLDivElement} task_list_node */
function init(task_list_node) {
    if (task_list_node.getAttribute('data-mounted') === 'true') {
        throw AlreadyMountedError(task_list_node)
    }

    /** @type {HTMLButtonElement[]} */
    const delete_task_btns = task_list_node.querySelectorAll(selectors.delete_task_btn);

    delete_task_btns.forEach((btn) => {
        btn.addEventListener('click', handleDeleteTaskClick);
        btn.setAttribute('data-mounted', 'true')
    })

    /** @type {HTMLButtonElement[]} */
    const view_task_btns = task_list_node.querySelectorAll(selectors.view_task_btn);

    view_task_btns.forEach((btn) => {
        btn.addEventListener('click', handleViewTaskClick);
        btn.setAttribute('data-mounted', 'true')
    })

    task_list_node.setAttribute('data-mounted', 'true')
}

/** @param {MouseEvent} ev */
async function handleViewTaskClick(ev) {
    /** @type {HTMLButtonElement} */
    const btn = ev.currentTarget;
    const task_id = btn.getAttribute('data-task-template-id');

    location.href = `/task_templates/${task_id}`
}

/** @param {MouseEvent} ev */
async function handleDeleteTaskClick(ev) {
    /** @type {HTMLButtonElement} */
    const btn = ev.currentTarget;
    const task_id = btn.getAttribute('data-task-template-id');

    const error = await service.deleteTask(task_id)

    if (error) {
        // TODO: parse error to user error
        console.error(error)
        return
    }

    document.querySelector(`[data-task-template-id="${task_id}"]`).remove()
}
