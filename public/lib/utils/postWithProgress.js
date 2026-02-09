/**
 * @typedef {Object} UploadResponse
 * @property {boolean} ok
 * @property {number} status
 * @property {string} statusText
 * @property {string} body
 */

/**
 * @param {string} url
 * @param {XMLHttpRequestBodyInit} body
 * @param {(percentage: number) => void} onProgress
 *
 * @returns {Promise<UploadResponse>}
 */
export default function postWithProgress(url, body, onProgress) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', url, true);

        let prev_progress = 0

        xhr.upload.addEventListener('progress', (event) => {
            if (event.lengthComputable) {
                const new_progress = Math.floor((event.loaded / event.total) * 100);

                if (new_progress > prev_progress) {
                    prev_progress = new_progress
                    onProgress(new_progress)
                }
            }
        });

        xhr.addEventListener('load', () => {
            resolve({
                ok: xhr.status >= 200 && xhr.status < 300,
                status: xhr.status,
                statusText: xhr.statusText,
                body: xhr.response,
            })
        })

        xhr.addEventListener('error', () => reject(new Error('Network error')));

        xhr.send(body);
    });
}
