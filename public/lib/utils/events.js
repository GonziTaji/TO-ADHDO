const events = Object.freeze({ register, dispatch })
export default events

export const EVENT_NAMES = /** @type {const} */ Object.freeze({
    new_task_template: "new:task_template",
})

/** @typedef {typeof EVENT_NAMES[keyof EVENT_NAMES]} CustomEventName */

/**
 * @typedef {(ev: CustomEvent<T>) => void} CustomEventHandler<T>
 * @template {{}} T
 * */

/** @type {{ event_name: CustomEventName, handler: CustomEventHandler }[]} */
const listeners = []

/**
 * @param {HTMLElement} source used to remove the listener if the source element no longer exists when the handler is called
 * @param {CustomEventName} event_name
 * @param {CustomEventHandler} handler
 */
function register(source, event_name, handler) {
    const handler_index = listeners.length

    /** @type {CustomEventHandler} */
    const autoRemoveHandler = (ev) => {
        if (!document.body.contains(source)) {
            const listener = listeners[handler_index]
            document.removeEventListener(listener.event_name, listener.handler)
            return;
        }

        handler(ev)
    }

    document.addEventListener(event_name, autoRemoveHandler)
    listeners.push({ event_name, handler: autoRemoveHandler })
}

/**
 * @param {CustomEventName} event_name
 * @param {CustomEvent<T>['detail']} payload
 * @template {{}} T
 */
function dispatch(event_name, payload) {
    console.log("DISPATCHING", event_name, payload)
    document.dispatchEvent(new CustomEvent(event_name, {
        detail: payload
    }))
}


