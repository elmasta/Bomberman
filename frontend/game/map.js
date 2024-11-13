const tileSize = 40,
    mapSizeX = 17,
    mapSizeY = 11;

let gameScreen

export const CreateGameScreen = () => {
    const gameScreenContainer = document.createElement("div")
    gameScreenContainer.style.height = (tileSize * mapSizeY) + 10 + "px"
    gameScreenContainer.style.width = (tileSize * mapSizeX) + 10  + "px"

    gameScreen = document.createElement("div")
    gameScreen.id = "gameScreen"
    gameScreen.tabIndex = 0;
    gameScreen.style.height = (tileSize * mapSizeY) + "px"
    gameScreen.style.width = (tileSize * mapSizeX) + "px"

    gameScreenContainer.append(gameScreen)

    return gameScreenContainer
}

export const CreateHUD = (player, nick) => {
    console.log(player)
    const HUD = document.createElement("div")
    HUD.id = player["User"]["Nickname"] + "-container"
    HUD.classList.add("playerInfo")
    switch (player["Number"]) {
        case 1:
            HUD.style.background = "yellow"
            break
        case 2:
            HUD.style.background = "blue"
            break
        case 3:
            HUD.style.background = "red"
            break
        case 4:
            HUD.style.background = "lime"
            break
    }

    if(nick == player["User"]["Nickname"]){
        HUD.classList.add("mainPlayer")
    }

    let nameHUD = document.createElement("div")
    nameHUD.classList.add("nameHUD")
    nameHUD.innerHTML = player["User"]["Nickname"]
    HUD.append(nameHUD)

    let innerHUD = document.createElement("div")
    for (let i=0; i<3; i++){
        let life = document.createElement("div")
        life.classList.add("life")
        innerHUD.appendChild(life)
    }    

    innerHUD.style.marginLeft = "10px"
    innerHUD.id = player["User"]["Nickname"] + "-HUDlives"
    innerHUD.classList.add("row")
    innerHUD.dataset.p_id = player["Number"]
    HUD.append(innerHUD)

    return HUD
}

export const UpdateTile = (posX, posY, type) => {
    if (type === ".") {
        console.log(posX, posY)
        RemoveTile(posX, posY)
    } else {
        let tile = GetTile(posX, posY)
        // get old type
        if (tile.classList[1] == "E"){
            SetExplosion(posX, posY)
        }
        if (tile.classList[1] == "R" || tile.classList[1] == "L" || tile.classList[1] == "M"){
            RemovePowerUp(posX, posY)
        }
        if (tile.classList[1] == "B"){
            BreakBrick(posX, posY)
        }
        //set new type
        tile.classList = "tile " + type

        if(type == "R" || type == "L" || type == "M"){
            SetPowerUp(posX, posY)
        }
    }
}

export const InitiateMap = (mapData) => {
    let posX = 0, posY = 0

    mapData.forEach((line) => {
        posX = 0
        line.forEach((tileType) => {
            if (tileType != ".") {
                if (isNaN(parseInt(tileType))) {
                    UpdateTile(posX, posY, tileType)
                } else {
                    CreatePlayer(tileType, posX, posY)
                }
            }
            posX++
        })
        posY++
    })
}

export const MovePlayer = (playerId, posX, posY) => {
    const scale = 1
    const player = GetPlayer(playerId);
    const translatedX = posX * tileSize;
    const translatedY = posY * tileSize;

    player.style.transform = `translate(${translatedX}px, ${translatedY}px) scale(${scale}, ${scale})`;
    player.style.transformOrigin = `${translatedX}px ${translatedY}px`;
    player.style.transformOrigin =  posX*tileSize + 'px ' + posY * tileSize + 'px'
}

export const UpdatePlayerHealth = (nbLives, operation, nickname) => {
    const playerHUD = document.getElementById(nickname + "-container")
    const playerLives = document.getElementById(nickname + "-HUDlives")
    if(nbLives == 0){
        playerHUD.style.background = "lightgray"
        playerHUD.style.opacity = "0.1"
    }

    // Animation
    if (operation == 0){
        playerLives.lastChild.classList.add("break")
        playerHUD.classList.add("shaking")
        const player = GetPlayer(playerLives.dataset.p_id)

        let span = document.createElement('span'); 
        span.style.backgroundImage = `url(../assets/P_${playerLives.dataset.p_id}_hurt.png)`;
        span.style.display = "flex";
        span.style.backgroundSize = "cover"
        span.style.backgroundRepeat = "no-repeat"
        span.classList.add("hurt");
        span.style.height = "40px";
        span.style.width = "40px";     
        player.appendChild(span);   
        player.classList.remove("p"+playerLives.dataset.p_id)    


         setTimeout(function(){
             player.classList.add("p"+playerLives.dataset.p_id)
             playerHUD.classList.remove("shaking")
             playerLives.lastChild.remove()

             span.remove()
          }, 1100)
    } else {
        let life = document.createElement("div")
        life.classList.add("life")
        playerLives.appendChild(life)
    }
}

const GetTile = (x, y) => {
    const tile = document.getElementById("tile-" + x + "-" + y)
    if (tile) { 
        return tile
    } else {
        return CreateTile(x, y)
    }
}

const GetPlayer = (playerId) => {
    return document.getElementById("p" + playerId)
}

const RemoveTile = (x, y) => {
    const tile = document.getElementById("tile-" + x + "-" + y)
    const pwrtile = document.getElementById("tile-" + x + "-" + y + "_pwr")
    if (tile) {
        tile.remove()
    }
    if (pwrtile) {
        pwrtile.remove()
    }
}

export const CreateTile = (posX, posY) => {
    const div = document.createElement('div')
    div.classList = "tile"
    div.id = "tile-" + posX + "-" + posY
    div.style.left = posX * tileSize + "px"
    div.style.top = posY * tileSize + "px"
    div.style.width = tileSize + "px"
    div.style.height = tileSize + "px"

    gameScreen.append(div)

    return div
}

export const CreatePlayer = (playerId, posX, posY) => {
    const div = document.createElement('div')
    div.classList = "player"        
    div.id = "p" + playerId
    div.classList.add(div.id) 
    div.style.left = "0px"
    div.style.top = "0px"
    div.style.width = tileSize + "px"
    div.style.height = tileSize + "px"
    div.style.transform = 'translate(' + posX*tileSize + 'px, ' + posY * tileSize + 'px)'

    gameScreen.append(div)
    
    return div
}

function SetExplosion(x, y){
    let blast = CreateTile(x, y)
    blast.id = ""
    blast.classList.add("blast")
    gameScreen.classList.add("shaking")

    setTimeout(function(){
        gameScreen.classList.remove("shaking")
    }, 500)
    setTimeout(function(){
        blast.remove()
    }, 1100)
}

function SetPowerUp(x, y){
    let halo = CreateTile(x, y)
    halo.id += "_pwr"
    halo.classList.add("Pwr")
}

function RemovePowerUp(x, y){
    let halo = document.getElementById("tile-"+x+"-"+y+"_pwr")
    halo.remove()
}

function BreakBrick(x, y){
    let b_broken = CreateTile(x, y)
    b_broken.id = ""
    b_broken.classList.add("B_broken") 

    setTimeout(function(){
        b_broken.remove()
    }, 500)

}

export function DisplayVictoryScreen(winner) {
    document.getElementById("victoryScreen").style.display = "flex"
   let victorytext = document.getElementById("victoryText")
  
    gameScreen.style.opacity = "0.3"
    if(winner == "p0"){
        document.getElementById("winner").style.display = "none"
        document.getElementById("trophy").style.display = "none"
        let displayLosers = document.getElementById("displaywinner")
        let i = 0 
        Array.from(document.getElementsByClassName("playerInfo")).forEach(elem => {
            console.log(elem)
            i++
            let loser = document.createElement("div")
            loser.classList.add(`p${i}`, `hurt`)
            loser.style.height = "100px"
            loser.style.width = "100px"
            loser.style.backgroundSize = "cover"
            displayLosers.appendChild(loser)
        })
            
        
        victorytext.innerHTML = "Nobody win !! "
    } else {        
        document.getElementById("winner").classList.add(winner)        
        victorytext.innerHTML += " Player " + winner[1]
        console.log("🏆", winner)
    }

    
}