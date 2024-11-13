import { ChangeElement } from "../framework/main.js"
// import { StartAnimating,  Stop} from "../game/fps.js"
import { CreateGameScreen, UpdateTile, UpdatePlayerHealth, MovePlayer, InitiateMap, CreateTile, CreatePlayer, CreateHUD, DisplayVictoryScreen } from "../game/map.js"
import { SetupInputListener } from "../game/playerInput.js"

//socket handling section
let socket = new WebSocket("ws://localhost:8080/ws")
let connected = false
let nick = ""

console.log("attempting web socket connection")
let test = document.getElementById("test")
//console.log( test.classList[1])

const CreateChatWindow = () => {
    const chatContainer = document.createElement("div")
    chatContainer.innerHTML = '' + 
    '<form id="chat-form">' +
    '<input type="text" placeholder="Enter message" >' +
    '</form>' + 
    '<div class="div" id="chat-window"></div>'
    
    chatContainer.id = "chat-container"
    
    return chatContainer
}

const CreateIntroForm = () => {
    const formContainer = document.createElement("div")
    formContainer.id = "introFormContainer"
    formContainer.innerHTML = '' + 
    '<form style="width: 100%;" id="intro-form">' +
    '<input type="text" placeholder="enter your nickname">' + 
    '</form>'

    return formContainer
}

const AddMessage = (message, origin) => {
    const div = document.createElement("div")
    div.innerText = origin +" : "+ message
    document.getElementById("chat-window").append(div)
}

const Initiate = () => {
    document.getElementById("introFormContainer")?.remove()
    document.body.append(CreateChatWindow())
    let gameContainer = document.getElementById("game-container")
    gameContainer.append(CreateGameScreen())
    var text = CreateTile(6, 3)
    text.classList.add("infotext")
    text.innerHTML = "Waiting for players... <span class=\"numbers\">"
    SetupInputListener(socket)
}

socket.onopen = () => {
    document.body.append(CreateIntroForm())
}

socket.onmessage = (event) => {
    var jsonData = JSON.parse(event.data)
    console.log(jsonData.datatype, jsonData.data )
    switch (jsonData.datatype) {
        case "chat_message":
            AddMessage(jsonData.data, jsonData.origin)
            break
        case "waiting_room":
            DisplayWaitingRoom(jsonData.data, jsonData.nbPlayer)
            break
        case "game_starting":
            DisplayStartingRoom(jsonData.data, jsonData.nbPlayer)
            break
        case "game_grid":
            ClearTiles()
            InitiateMap(jsonData.data)
            let gameContainer = document.getElementById("game-container")
            let HUDContainer = document.createElement("div")
            HUDContainer.classList.add("HUDContainer")
            for (let i in jsonData.players) {
                CreatePlayer(jsonData.players[i]["Number"], jsonData.players[i]["Position"][1], jsonData.players[i]["Position"][0])
                HUDContainer.append(CreateHUD(jsonData.players[i], nick))
            }

            gameContainer.append(HUDContainer)
            // StartAnimating(62)
            break
        case "loop_result":
            if (jsonData.data == "endgame") {
                //end the game
                // Stop = true //break the FPS counter loop
                setTimeout(function(){
                    let winner = "p0"
                    let playersLeft = document.getElementsByClassName("player")
                    if (playersLeft.length > 0){
                        winner = playersLeft[0].classList[1]
                    }
                    DisplayVictoryScreen(winner)
                }, 600)
                
            } else {
                console.log(Array.isArray(jsonData.data))
                if (!Array.isArray(jsonData.data)){
                    UpdatePlayerHealth(jsonData.data.Y, jsonData.data.X, jsonData.data.Extra)
                }

                Array.from(jsonData.data).forEach(elem => {
                    if (elem.SquareType == "H") {
                        UpdatePlayerHealth(elem.Y, elem.X, elem.Extra)
                    }
                    if (elem.SquareType == "D") {
                        setTimeout(function(){
                            const user = document.getElementById("p" + elem.Extra)
                            console.log("dead", user)
                            user.remove()
                        }, 500)
                        
                    } else {
                        UpdateTile(elem.X, elem.Y, elem.SquareType)
                    }
                })
            }
            break
        case "player_move":
            console.log(jsonData.data)
            MovePlayer(jsonData.data.PlayerNMB, jsonData.data.PlayerPosition[1], jsonData.data.PlayerPosition[0])
            break
        case "remove_powerup":
            UpdateTile(jsonData.data[1], jsonData.data[0], ".")
            break
        case "bomb":
            UpdateTile(jsonData.data[1], jsonData.data[0], "E")
            break
    }
}

socket.onclose = (event) => {
    document.body.innerHTML = ""
    document.body.append(CreateIntroForm())
}

document.body.addEventListener("submit", (event) => {
    event.preventDefault()
    switch (event.target.id) {
        case "intro-form":
            nick = event.target.children[0].value
            socket.send(JSON.stringify({
                datatype: "nick_setting",
                message: event.target.children[0].value
            }))
            Initiate()
            socket.send(JSON.stringify({
                datatype: "ready_to_play",
                message: true
            }))
            break
        case "chat-form":
            socket.send(JSON.stringify({
                datatype: "chat_message",
                message: event.target.children[0].value
            }))
            event.target.children[0].value = ""
            break
    }
})

function ClearTiles(){
    var tiles = document.getElementsByClassName("tile")
    Array.from(tiles).forEach(div =>
        div.remove()
    )
}

function DisplayWaitingRoom(time, nbPlayer){    
    ClearTiles()
    var text = CreateTile(6, 3)
    text.classList.add("infotext")
    text.innerHTML = "Waiting for players... <span class=\"numbers\">" + time + "</span>"
    for (let i=1; i<=4; i++){
        var tile = CreateTile(i+5, 4)
        if(nbPlayer>=i){
            tile.classList.add("p"+i) 
        } else {
            tile.classList.add("X")
            var loadingTile = CreateTile(i+5, 4)
            loadingTile.classList.add("Loading")
        }
    }
}

function DisplayStartingRoom(time, nbPlayer){    
    ClearTiles()
    var text = CreateTile(6, 3)
    text.classList.add("infotext")
    text.innerHTML = "Game is starting ! <span class=\"numbers\">" + time + "</span>"
    for (let i=1; i<=nbPlayer; i++){
        var tile = CreateTile(i+5, 4)
            tile.classList.add("p"+i) 
    }   
}




