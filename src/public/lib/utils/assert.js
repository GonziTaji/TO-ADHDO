/**
 * @param {unknown} test
 * @param {string} message
 *
 * @satisfies {typeof test}
 * */
export default function assert(test, message) {
    if (!test) {
        throw new Error(message)
    }
}
