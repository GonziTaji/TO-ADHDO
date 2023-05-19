export function validateParams(body: any, params: string[]) {
    for (const param of params) {
        if (!body[param]) {
            return 'Missing parameter: ' + param;
        }
    }

    return null;
}
