const controller = {
    mount,
}

export default controller

/** @type {HTMLDivElement} */
let content_container
/** @type {HTMLDivElement} */
let buttons_container
/** @type {string} */
let task_id

/**
 * @param {HTMLDivElement} param_task_id
 * @param {HTMLDivElement} param_content_container
 * @param {HTMLDivElement} param_buttons_container
 */
function mount(param_task_id, param_content_container, param_buttons_container) {
    task_id = param_task_id
    content_container = param_content_container
    buttons_container = param_buttons_container

    buttons_container.querySelectorAll('button[data-btn-action]').forEach((btn) => {
        const action = btn.getAttribute('data-btn-action')

        switch (action) {
            case 'edit':
                return btn.addEventListener('click', handleEditClick)
            case 'delete':
                return btn.addEventListener('click', handleRemoveClick)
            case 'save':
                return btn.addEventListener('click', handleSaveClick)
            case 'cancel':
                return btn.addEventListener('click', handleCancelClick)
        }
    })
}

function handleEditClick() {
    setForm()
}

function handleCancelClick() {
    setView()
}

async function handleRemoveClick() {
    const r = await fetch(`/api/task_templates/${task_id}`, { method: 'DELETE' })

    if (!r.ok) {
        alert(String(err))
        return
    }

    location.href = '/'
}

async function handleSaveClick() {
    const form_data = {}

    const r = await fetch(`/api/task_templates/${task_id}`, {
        method: 'PUT',
        body: JSON.stringify(form_data)
    })

    if (!r.ok) {
        alert(String(err))
        return
    }

    await setView()
}

async function setForm() {
    const r = await fetch(`/api/task_templates/${task_id}/form`)

    if (!r.ok) {
        alert(r.statusText)
        return
    }

    buttons_container.setAttribute('data-mode', 'edit')

    content_container.innerHTML = await r.text()

    content_container.querySelector('form').addEventListener('submit', () => {
        setView()
    })
}

async function setView() {
    const r = await fetch(`/api/task_templates/${task_id}/view`)

    if (!r.ok) {
        alert(r.statusText)
        return
    }

    content_container.innerHTML = await r.text()
    buttons_container.setAttribute('data-mode', 'normal')
}
