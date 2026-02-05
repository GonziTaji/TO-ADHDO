/**
 * @param {string} url
 * @param {XMLHttpRequestBodyInit} body
 * @param {(percentage: number) => void} onProgress
 *
 * @returns {Promise<{ ok: boolean, status: number, statusText: string, body: string }>}
 */
export default function postWithProgress(url, body, onProgress) {
    return new Promise((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.open('POST', url, true);

        let prev_progress = 0

        xhr.upload.addEventListener('progress', (event) => {
            if (event.lengthComputable) {
                const new_progress = Math.floor((event.loaded / event.total) * 100);
                console.log(`Upload progress: ${new_progress}%`);

                if (new_progress > prev_progress) {
                    prev_progress = new_progress
                    onProgress(new_progress)
                }
            }
        });

        xhr.addEventListener('load', () => {
            if (xhr.status >= 200 && xhr.status < 500) {
                resolve({
                    ok: xhr.status == 200,
                    status: xhr.status,
                    statusText: xhr.status,
                    body: xhr.response
                });
            } else {
                reject(new Error(xhr.statusText));
            }
        });

        xhr.addEventListener('error', () => reject(new Error('Network error')));

        xhr.send(body);
    });
}
