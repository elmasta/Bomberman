## Server to client messages, datatype 'game-event'

List of all events that can be sent to clients from the server.

### update_tile(type, coord)

Tile at 'coord' changes state to 'type'. New type can be :
- '.' : empty tile (bloc destroyed, bonus picked-up, flame stopped)
- 'r', 'l', 'm' : tile with bonus that can be picked-up (looted by destroyed bloc)
- 'f' : flame (bomb explosion)
- 'b' : bomb (placed by player)

### player_move(playedId, coord)

Player 'playerId' moving to tile at 'coord'. 

### player_touched(playerId)

Player 'playerId' lost a life. He should flash.

### player_dead(playerId)

Player 'playerId' is dead. He should be removed from the screen.

### game_end(playerId)

Game is finished, won by player 'playerId'.