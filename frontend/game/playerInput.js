const PressedKeys = {}

export const SetupInputListener = (ws) => {
    document.getElementById("gameScreen").focus()

    document.getElementById("gameScreen").addEventListener("click", (e) => {
        e.target.focus()
    })

    document.body.addEventListener("keydown", (e) => {
        if (document.activeElement.id != "gameScreen") {
            return
        }
        const key = ParseKey(e)
        if (key != undefined && !PressedKeys[key]) {
            PressedKeys[key] = true
            ws.send(JSON.stringify({datatype: "key_press", data: {key: key, status: "on"}}))       
        }
    })

    document.body.addEventListener("keyup", (e) => {
        if (document.activeElement.id != "gameScreen") {
            return
        }
        const key = ParseKey(e)
        if (key) {
            PressedKeys[key] = false
            ws.send(JSON.stringify({datatype: "key_press", data: {key: key, status: "off"}}))
        }
    })
}

const ParseKey = (e) => {
    let key;
    switch (e.key) {
    case 'ArrowUp':
        key = "up"
        break;
    case 'ArrowDown': 
        key = "down"
        break;
    case 'ArrowLeft':  
        key = "left"
        break;
    case 'ArrowRight': 
        key = "right"
        break;
    case ' ':
        key = "bomb"
        break;
    }

    return key
}