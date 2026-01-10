export default function onElementRemoved(element, callback) {
    if (!element) return () => { };

    let called = false;

    const observer = new MutationObserver(() => {
        if (!document.body.contains(element)) {
            if (called) return;
            called = true;
            observer.disconnect();
            callback();
        }
    });

    observer.observe(document.body, {
        childList: true,
        subtree: true
    });

    return () => observer.disconnect();
}

