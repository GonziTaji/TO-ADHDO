
/**
 * A script wanted to control an element that is already being controlled (controlled already mounted)
 *
 * @param {HTMLElement} node
 */
export const AlreadyMountedError = (node) => CustomError('AlreadyMountedError', `Element already mounted:\n${node.outerHTML}`)

/**
 * An element that was expected to have tagName = `expected`, has a different tagName
 *
 * @param {string} expected
 * @param {string} got
 */
export const UnexpectedTagNameError = (expected, got) =>
    CustomError(
        'UnexpectedTagNameError',
        `Element ${node.tagName} was expected as "${expected}". Got "${got}"`
    )

/**
 * An https call to a server failed with code `code`
 *
 * @typedef {Error & { code: string }} HttpServiceError 
 *
 * @param {string} code
 * @param {string?} message
 *
 * @returns {HttpServiceError}
 * */
export const HttpServiceError = (code, message = "") => {
    /** @type {HttpServiceError} */
    const error = CustomError('HttpServiceError', message)
    error.code = code

    return error
}

const errors = {
    AlreadyMountedError,
    UnexpectedTagNameError,
}

export default errors

/**
 * @param {string} type_name
 * @param {unknown} payload
 */
function CustomError(type_name, payload) {
    const error = new Error(String(payload))
    error.name = type_name
    return error
}
