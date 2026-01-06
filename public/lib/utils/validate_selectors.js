
/**
 * @typedef {} 
 */


/**
 * Validates that each selector in `selectors` matches with an existing element inside the parent
 *
 * @param {HTMLElement} parent
 * @param {string[]} selectors
 *
 * @returns {Error | null} error
 * */
export default function validateSelectors(parent, selectors) {
    const selectors_not_found = selectors
        .filter((selector) => parent.querySelector(selector) === null)

    if (selectors_not_found.length > 0) {
        const msg = 'Missing elements in provided parent:' + selectors_not_found.join(', ')
        return new Error(msg)
    }

    return null
}

