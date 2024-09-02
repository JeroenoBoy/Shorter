const observer = new MutationObserver((records) => {
    records.forEach(record => {
        record.addedNodes.forEach(node => {
            setTimeout(() => {
                node.remove()
            }, 10_000)
        })
    })
})

const notifications = document.querySelector("#notifications-container")
observer.observe(notifications, {
    childList: true,
    subtree: false,
    attributes: false,
    characterData: false
})
