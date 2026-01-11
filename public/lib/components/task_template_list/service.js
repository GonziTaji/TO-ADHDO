import { HttpServiceError } from "../../utils/errors.js";

const service = {
    getTaskListItem,
    deleteTask,
}

export default service

/**
 * @param {string} task_id
 * @returns {Promise<HttpServiceError | null>}
 */
async function deleteTask(task_id) {
    const res = await fetch(`/api/task_templates/${task_id}`, { method: "DELETE" })

    if (!res.ok) {
        console.error("Failed to delete task_template", res.status, res.statusText)

        return HttpServiceError(res.status, "Task template could not be deleted")
    }

    return null
}

/**
 * @param {string} task_id
 * @returns {Promise<HttpServiceError | string>}
 */
async function getTaskListItem(task_id) {
    const res = await fetch(`/api/task_templates/${task_id}/list-item`)

    if (!res.ok) {
        console.error("Failed to get task_template as list item", res.status, res.statusText)

        return HttpServiceError(res.status, "Task template could not be loaded")
    }

    return res.text()
}
