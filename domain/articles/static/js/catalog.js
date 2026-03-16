function getFiltersForm() {
    return document.getElementById('filters-form')
}

async function renderList() {
    const fd = new FormData(getFiltersForm())
    const sp = new URLSearchParams(fd)

    const history_pathname = `${location.pathname}?${sp.toString()}`
    const history_url = new URL(history_pathname, location.origin)
    history.pushState({}, '', history_url.toString())

    const fetch_pathname = `${location.pathname}list?${sp.toString()}`
    const fetch_url = new URL(fetch_pathname, location.origin)

    try {
        const res = await fetch(fetch_url.toString())
        const html = await res.text()

        document.getElementById('catalog-list').outerHTML = html
    } catch (e) {
        console.warn('could not render list', e)
        alert("Algo paso! Recarga la pagina para intentarlo de nuevo")
    }
}

document.addEventListener('DOMContentLoaded', () => {
    const filtersForm = getFiltersForm()

    filtersForm.addEventListener('change', (ev) => {
        const tagInput = ev.target.closest('[name="tags"]')

        if (tagInput) {
            renderList()
        }
    })

    filtersForm.addEventListener('submit', (ev) => {
        ev.preventDefault()
        renderList()
    })

    filtersForm.addEventListener('reset', () => {
        renderList()
    })
})

document.addEventListener('reset', (ev) => {
    ev.preventDefault()

    const form = ev.target.closest('form')

    if (!form) return

    switch (form.id) {
        case "search-form": {
            const url = new URL(location)
            url.searchParams.delete("s")
            location = url
            break;
        }

        case "filters-form": {
            const url = new URL(location)
            url.searchParams.delete('tags')
            location = url
            break;
        }
    }
}, { capture: true })

