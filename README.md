# TO-ADHDO

## Description

This app pivoted from a glorified todo-list into a glorified shop or something lol.

The goal of this project is to practice golang and reinforce web technologies knowledge.

## Objectives

- Reinforce knowledge of native web technologies
- Become more familiar with "new" (modern) features of web apis, css, js, etc.
- Delve deeper into golang's standart library
- Discover common patterns for a golang web server, and for webservers in general

## Expected result

A scalable and maintanable proof of concept of a competent webapp without frontend frameworks

## Restrictions

To avoid losing direction of this project's goals, some dev restrictions are defined:

### DO

- Compose the web pages using html/templates' block and template features
- Treat hypermedia as the application state
- Use go packages 

### DON'T

- Import javascript libraries unless strictly necessary
- Install golang packages besides the obv . Stay away from "frameworks" and leaky abstractions
- 

### AVOID

- heavy DOM manipulation. Prefer manipulating the DOM with js scripts.

### ONLY

- ONLY use the standart library and go-gin for the server (exceptions may exist)
- ONLY USE hypermedia as the frontend's state (htmx's philosophy)
- ONLY send html (or text for info) from the server to the client, unless strictly necessary


Examples (frontend):

Instead of updating DOM elements with js:
- Request the component/section/slice of html and replace the stale elements

Instead of creating or duplicating html nodes with js:
- Use native html template elements

Instead of using complex css selectors
- Use modern css features like nested selectors, and pseudoclasses like :has

Examples (backend):

idk, that's why I'm doing this lol

## Tech Stack

- The web pages are a composition of html/template blocks written in plain html/css/js
- The backend is a go server written with go-gin
- The data is saved in an SQLite instance living where the server lives (for now I guess)

