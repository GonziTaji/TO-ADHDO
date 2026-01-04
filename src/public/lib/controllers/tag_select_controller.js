import elements from "../utils/elements.js";

document.addEventListener("DOMContentLoaded", init);

function init() {
    // events
    const input = elements.task_name_input();
    input.addEventListener("input", handleTagInputChange);
    input.addEventListener("keydown", handleTagInputKeyDown);
}

/** @param {KeyboardEvent} ev */
function handleTagInputChange(ev) {
    const option_tags_nodes = [
        ...elements.tags_options_list().querySelectorAll("option"),
    ];

    const search = ev.currentTarget.value.toLowerCase().trim();

    const filtered_tags = !search ? [] : option_tags_nodes
        .filter((option) => !option.disabled)
        .filter((option) => option.value.includes(search))
        .map((option) => option.value);

    const ul = elements.filtered_tags_list();

    ul.innerHTML = "";

    for (const tag of filtered_tags) {
        const btn = document.createElement("button");

        btn.type = "button";
        btn.innerText = tag;

        btn.addEventListener("click", () => {
            addTagToTask(tag);
        });

        const li = document.createElement("li");

        li.appendChild(btn);
        ul.appendChild(li);
    }
}

/** @param {KeyboardEvent} ev */
function handleTagInputKeyDown(ev) {
    const value = ev.currentTarget.value.trim().toLowerCase();

    if (value === "") {
        return;
    }

    if (ev.key === "Enter" || ev.key === "Tab") {
        ev.preventDefault();
        addTagToTask(value);
    }
}

/** @param {string} tag_name */
function addTagToTask(tag_name) {
    /** @type {HTMLTemplateElement} */
    const node = elements.selected_tag_template().content.cloneNode(true);

    const tag_card = node.querySelector(".selected_tag_card");

    const tag_option = elements
        .tags_options_list()
        .querySelector(`option[value="${tag_name}"]`);

    tag_card.querySelector("input").value = tag_name;
    tag_card.querySelector("button").addEventListener("click", () => {
        tag_card.remove();

        if (tag_option) {
            tag_option.disabled = false;
        }
    });

    elements.selected_tags_list().appendChild(tag_card);
    elements.task_name_input().value = "";
    elements.filtered_tags_list().innerHTML = "";

    if (tag_option) {
        tag_option.disabled = true;
    }
}
