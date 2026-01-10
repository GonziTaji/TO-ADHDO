import { AlreadyMountedError } from "../../utils/errors.js";
import { EVENT_NAMES } from "../../utils/events.js";
import onElementRemoved from "../../utils/on_element_removed.js";
import service from "./service.js";

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

    registerButtonEvents(task_list_node)

    const handler = (ev) => handleTaskCreated(task_list_node, ev)

    document.addEventListener(EVENT_NAMES.new_task_template, handler)

    onElementRemoved(
        task_list_node,
        () => document.removeEventListener(EVENT_NAMES.new_task_template)
    )

    task_list_node.setAttribute('data-mounted', 'true')
}

/**
 * @param {HTMLDivElement} task_list_node
 */
function registerButtonEvents(task_list_node) {
    /** @type {HTMLButtonElement[]} */
    const delete_task_btns = task_list_node.querySelectorAll(
        `${selectors.delete_task_btn}:not([data-mounted="true"])`
    );

    delete_task_btns.forEach((btn) => {
        btn.addEventListener('click', handleDeleteTaskClick);
        btn.setAttribute('data-mounted', 'true')
    })

    /** @type {HTMLButtonElement[]} */
    const view_task_btns = task_list_node.querySelectorAll(
        `${selectors.view_task_btn}:not([data-mounted="true"])`
    );

    view_task_btns.forEach((btn) => {
        btn.addEventListener('click', handleViewTaskClick);
        btn.setAttribute('data-mounted', 'true')
    })
}

/**
 * @param {HTMLDivElement} task_list_node
 * @param {CustomEvent<{id:string}>} ev
 */
async function handleTaskCreated(task_list_node, ev) {
    console.log("handleTaskCreated called")
    const task_id = ev.detail.id

    const response = await service.getTaskListItem(task_id)

    console.log("response for task list item id: " + task_id, response)

    if (typeof response !== 'string') {
        alert(response.message)
        return
    }

    const ul = task_list_node.querySelector('ul')

    ul.insertAdjacentHTML('afterbegin', response)
    const li = ul.firstChild

    console.log("new li", li)

    registerButtonEvents(task_list_node)
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
