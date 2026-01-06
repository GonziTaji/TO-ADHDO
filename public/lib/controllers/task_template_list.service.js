import { HttpServiceError } from "../utils/errors.js";

const service = {
    deleteTask,
}

export default service

/**
 * @param {string} task_id
 * @returns {Promise<HttpServiceError | null>}
 */
async function deleteTask(task_id) {
    const res = await fetch(`/api/tasks_templates/${task_id}`, { method: "DELETE" })

    if (!res.ok) {
        console.error("Failed to delete task_template", res.status, res.statusText)

        return HttpServiceError(res.status, "Task template could not be deleted")
    }

    return null
}
