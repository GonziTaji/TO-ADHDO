/**
 * @param {HTMLTemplateElement} template
 * @returns {HTMLElement}
 * */
export function getFirstChildCopyFromTemplate(template) {
    const child = template.content.cloneNode(true).firstElementChild;
    return child
}
