import assert from "./assert.js"

const elements = {
    /** @returns {HTMLDataListElement} */
    tags_options_list: () => document.getElementById("datalist_available_tags"),
    /** @returns {HTMLInputElement} */
    task_name_input: () => document.getElementById("input_create_task_tag_selector"),
    /** @returns {HTMLTemplateElement} */
    selected_tag_template: () => document.getElementById("template_create_task_selected_tag"),
    /** @returns {HTMLUListElement} */
    filtered_tags_list: () => document.getElementById("ul_filtered_tags"),
    /** @returns {HTMLDivElement} */
    selected_tags_list: () => document.getElementById("div_selected_tags"),
    /** @returns {HTMLFormElement} */
    task_form: () => document.getElementById("task_form"),
    /** @returns {HTMLTextAreaElement} */
    task_description_textarea: () => document.getElementById("textarea_create_task_description"),
    /** @returns {HTMLDivElement} */
    task_list_container: () => document.getElementById("task_list_container"),
};

export default elements

document.addEventListener("DOMContentLoaded", assertElements);

function assertElements() {
    Object.entries(elements).forEach(([key, fn]) =>
        assert(fn(), `task form element of key "${key}" not found in DOM`)
    )
}
