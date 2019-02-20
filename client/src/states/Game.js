/* globals __DEV__ */
import Phaser from 'phaser';
import Player from '../sprites/Player';
import * as effect from '../sprites/Effect';
import Map from '../sprites/Map';
import Leaderboard from '../sprites/Leaderboard';
import messagepb from './message_pb';
import config from '../config';

// __HOST_IP__ is server IP
var host = 'ws:///' + __HOST_IP__ + '/game/';
console.log('Connecting to ' + host);

const SERVER_MESSAGE_CASES = messagepb.ServerGameMessage.MessageCase;
const CLIENT_MESSAGE_CASES = messagepb.ClientGameMessage.MessageCase;

export default class extends Phaser.State {
  init () {
    this.socket = new WebSocket(host);
    this.socket.binaryType = 'arraybuffer';

    // list of players
    this.playerList = [];
    // layerOrder is display ordering, the prior is after the latter
    this.layerOrder = [];
  }
  preload () {}

  create () {
    this.objectLayer = this.game.add.group();
    this.backgroundLayer = this.game.add.group();
    this.uiLayer = this.game.add.group();
    this.layerOrder = [this.backgroundLayer, this.objectLayer, this.uiLayer];
    //this.glowFilter = new Phaser.Filter.Glow(this.game);
    //this.shaders = [ this.glowFilter ];
    // No filter for better performance
    this.shaders = [];

    //  Our tiled scrolling background
    this.game.scale.pageAlignHorizontally = true;
    this.game.scale.pageAlignVertically = true;
    this.game.scale.refresh();

    // Remove background
    //this.land = this.game.add.tileSprite(0, 0, config.screenWidth, config.screenHeight, 'maptile');
    //this.land.fixedToCamera = true;
    //this.backgroundLayer.add(this.land);
    this.game.stage.backgroundColor = '#000000';

    // Game world bound
    this.game.world.setBounds(0, 0, 19200, 19200);

    // Create leaderboard
    this.leaderboard = new Leaderboard({
      game: this.game,
      x: 100,
      y: 10,
      width: 30,
      height: 100,
    });
    this.leaderboard.fixedToCamera = true;
    this.uiLayer.add(this.leaderboard);

    this.cursors = this.game.input.keyboard.createCursorKeys();
    this.wasd = {
      left: this.game.input.keyboard.addKey(Phaser.Keyboard.A),
      right: this.game.input.keyboard.addKey(Phaser.Keyboard.D),
      up: this.game.input.keyboard.addKey(Phaser.Keyboard.W),
      down: this.game.input.keyboard.addKey(Phaser.Keyboard.S)
    };

    this.cntSeqNum = 0;
    this.pending_inputs = [];
    var emitter = this.game.add.emitter(200, 200, 200);

    this.uiLayer.add(emitter);
    this.setupEvent(this.socket, this);
  }

  render () {
    if (__DEV__) {
      // this.game.debug.spriteInfo(this.player, 32, 32);
    }
  }

  createInitPlayerMessage (clientID, name) {
    var message = new messagepb.ClientGameMessage();
    var initPlayer = new messagepb.InitPlayer();
    initPlayer.setName(name);
    initPlayer.setClientId(clientID);
    message.setInitPlayerPayload(initPlayer);

    return message.serializeBinary();
  }

  createMoveProtoMessage (id, timeElapsed, dx, dy) {
    this.cntSeqNum++;
    var message = new messagepb.ClientGameMessage();
    var movePosition = new messagepb.MovePosition();
    movePosition.setId(id);
    movePosition.setDx(dx);
    movePosition.setDy(dy);
    message.setMovePositionPayload(movePosition);
    message.setTimeElapsed(timeElapsed);
    message.setInputSequenceNumber(this.cntSeqNum);

    return message.serializeBinary();
  }

  createShootProtoMessage (playerid, timeElapsed, camX, camY, x, y, mouseX, mouseY) {
    this.cntSeqNum++;
    var message = new messagepb.ClientGameMessage();
    var shoot = new messagepb.Shoot();
    var length = Math.sqrt(Math.pow(mouseX - (x - camX), 2) + Math.pow(mouseY - (y - camY), 2));
    shoot.setPlayerId(playerid);
    shoot.setX(x);
    shoot.setY(y);
    shoot.setDx((mouseX - (x - camX)) / length);
    shoot.setDy((mouseY - (y - camY)) / length);
    message.setShootPayload(shoot);
    message.setTimeElapsed(timeElapsed);
    message.setInputSequenceNumber(this.cntSeqNum);

    return message.serializeBinary();
  }

  update () {
    if (this.gameStart === undefined || this.gameStart === false) {
      return;
    }

    var timeNow = new Date();
    var lastTimestamp = this.lastTimestamp || timeNow;
    var timeElapsed = (timeNow - lastTimestamp) / 1000.0;
    var dx = 0;
    var dy = 0;
    this.lastTimestamp = timeNow;

    // Handle input
    // 1 is just to represent the direction
    if (this.cursors.left.isDown || this.wasd.left.isDown) {
      dx -= 1;
    }
    if (this.cursors.right.isDown || this.wasd.right.isDown) {
      dx += 1;
    }
    if (this.cursors.up.isDown || this.wasd.up.isDown) {
      dy -= 1;
    }
    if (this.cursors.down.isDown || this.wasd.down.isDown) {
      dy += 1;
    }

    if (dx !== 0 || dy !== 0) {
      var moveMessage;
      // Send move message to server
      moveMessage = this.createMoveProtoMessage(this.player.id, timeElapsed, dx, dy);
      this.socket.send(moveMessage);
      this.applyInput(moveMessage);
    }

    if (this.game.input.activePointer.leftButton.isDown) {
      // Send shoot messsage to server
      var shootMessage = this.createShootProtoMessage(this.player.id, timeElapsed, this.camera.x, this.camera.y, this.player.x, this.player.y, this.game.input.activePointer.x, this.game.input.activePointer.y);
      this.pending_inputs.push(shootMessage);
      this.socket.send(shootMessage);
    }

    // Check bullet hit block
    this.playerList.forEach(player => {
      // Get all bullets from all player
      player.shootManager.forEachAlive(bullet => {
        // Check if bullet is in Map
        if (!this.isInMap(bullet.position.x, bullet.position.y)) {
          effect.explode_bullet(this.game, this.uiLayer, bullet.position.x, bullet.position.y);
          bullet.kill();
        }
        // Check if bullet hit block
        if (this.isBulletHitBlock(bullet.position.x, bullet.position.y)) {
          effect.explode_bullet(this.game, this.uiLayer, bullet.position.x, bullet.position.y);
          bullet.kill();
        }

        // check if bullet hit enemies
        if (this.isBulletHitPlayers(player, this.playerList, bullet.position.x, bullet.position.y)) {
          effect.explode_bullet(this.game, this.uiLayer, bullet.position.x, bullet.position.y);
          effect.explode_hit(this.game, this.uiLayer, bullet.position.x, bullet.position.y);
          bullet.kill();
        }
      });
    });

    // Update scoreboard
    this.leaderboard.updateLeaderboard(this.playerList);

    // Display layer in order
    for (var i = 0; i < this.layerOrder.length; i++) {
      this.game.world.bringToTop(this.layerOrder[i]);
    }
  }

  isBulletHitPlayers (player, playerList, x, y) {
    for (var enemy of playerList) {
      if (enemy !== player && enemy.isCollidePoint(x, y)) {
        return true;
      }
    }
    return false;
  }

  isBulletHitBlock (x, y) {
    return this.map.isCollide(x, y);
  }

  isInMap (x, y) {
    return this.map.isInMap(x, y);
  }

  applyInput (data) {
    var msg = messagepb.ClientGameMessage.deserializeBinary(data);
    switch (msg.getMessageCase()) {
      case CLIENT_MESSAGE_CASES.MOVE_POSITION_PAYLOAD:
        msg = msg.getMovePositionPayload();
        var player = this.getPlayerByID(msg.getId());
        if (player === null) {
          return;
        }
        this.player.move(msg.getDx(), msg.getDy());
        break;
    }
  }

  bulletHitEnemy (tank, bullet) {
    effect.explode_bullet(this.game, this.uiLayer, bullet.position.x, bullet.position.y);
    bullet.kill();
  }

  setupEvent (socket, game) {
    socket.onopen = function (event) {
    };

    socket.onmessage = function (event) {
      var msg = messagepb.ServerGameMessage.deserializeBinary(event.data);
      switch (msg.getMessageCase()) {
        case SERVER_MESSAGE_CASES.UPDATE_PLAYER_PAYLOAD:
          game.updatePlayer(msg.getUpdatePlayerPayload());
          break;
        case SERVER_MESSAGE_CASES.INIT_SHOOT_PAYLOAD:
          game.initShoot(msg.getInitShootPayload());
          break;
        case SERVER_MESSAGE_CASES.INIT_ALL_PAYLOAD:
          game.initAll(msg.getInitAllPayload());
          break;
        case SERVER_MESSAGE_CASES.INIT_PLAYER_PAYLOAD:
          game.initPlayer(msg.getInitPlayerPayload());
          break;
        case SERVER_MESSAGE_CASES.REMOVE_PLAYER_PAYLOAD:
          game.removePlayer(msg.getRemovePlayerPayload());
          break;
        case SERVER_MESSAGE_CASES.REGISTER_CLIENT_ID_PAYLOAD:
          game.registerClientID(msg.getRegisterClientIdPayload());
          // Receive clientID, then register with name
          game.socket.send(
            game.createInitPlayerMessage(game.clientID, game.state.playerName)
          );
      }
    };
  }

  updatePlayer (playerMsg) {
    var player = this.getPlayerByID(playerMsg.getId());
    if (player === null) {
      return;
    }

    if (player.x !== playerMsg.getX() || player.y !== playerMsg.getY()) {
      player.x = playerMsg.getX();
      player.y = playerMsg.getY();
      // Move animation
      player.move();
    }
    player.health = playerMsg.getHealth();
    player.score = playerMsg.getScore();
  }

  getPlayerByID (id) {
    var filteredPlayers = this.playerList.filter(player => player.id === id);
    if (filteredPlayers.length === 0) {
      return null;
    }
    return filteredPlayers[0];
  }

  initPlayer (initPlayerMsg) {
    var player = new Player({
      game: this.game,
      layer: this.objectLayer,
      shaders: this.shaders,
      x: initPlayerMsg.getX(),
      y: initPlayerMsg.getY(),
      id: initPlayerMsg.getId(),
      name: initPlayerMsg.getName(),
      asset: 'tank'
    });
    this.playerList.push(player);

    if (initPlayerMsg.getIsMain()) {
      this.player = player;
      this.game.camera.follow(this.player);
    }
  }

  initAll (initAllMsg) {
    // This initALl package contain all the map info and players info
    var initMap = initAllMsg.getInitMap();
    var initPlayers = initAllMsg.getInitPlayerList();
    // Init map
    this.map = new Map({
      game: this.game,
      layer: this.backgroundLayer,
      shaders: this.shaders,
      blockWidth: initMap.getBlockWidth(),
      blockHeight: initMap.getBlockHeight(),
      numCols: initMap.getNumCols(),
      numRows: initMap.getNumRows(),
      blocks: initMap.getBlockList()
    });

    // Init current players in the scene
    initPlayers.forEach(
      player => {
        this.initPlayer(player);
      }
    );

    this.gameStart = true;
  }

  removePlayer (removePlayerMsg) {
    var playerID = removePlayerMsg.getId();
    var player = this.getPlayerByID(playerID);
    if (player == null) {
      return;
    }

    for (var i = 0; i < this.playerList.length; i++) {
      if (this.playerList[i].id === playerID) {
        this.playerList.splice(i, 1);
        break;
      }
      this.playerList.filter(player => player.id !== playerID);
    }

    // Player clean
    player.destroy();
    player.removeChildren(0, player.length);
    player.healthbar.destroy();
    player.nametag.destroy();
    player.shootManager.destroy();
    player.emitter.destroy();
    // Exploding effect
    effect.explode(this.game, this.uiLayer, this.shaders, player.x, player.y);
    this.objectLayer.remove(player);

    if (player === this.player) {
      // Rejoin
      alert('You are killed. Press OK to restart the game', this.socket.send(
        this.createInitPlayerMessage(this.clientID, this.player.name)));
    }
  }

  registerClientID (registerClientIDMsg) {
    this.clientID = registerClientIDMsg.getClientId();
  }

  initShoot (initShootMsg) {
    var playerID = initShootMsg.getPlayerId();
    var player = this.getPlayerByID(playerID);
    if (player === null) {
      return;
    }
    player.fire(
      initShootMsg.getX(),
      initShootMsg.getY(),
      initShootMsg.getDx(),
      initShootMsg.getDy()
    );
  }
}
