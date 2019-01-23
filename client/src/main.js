import 'pixi';
import Phaser from 'phaser';

import BootState from './states/Boot';
import SplashState from './states/Splash';
import GameState from './states/Game';

import config from './config';

class Game extends Phaser.Game {
  constructor () {
    // TODO: fullscreen but keep ratio
    super(config.screenWidth, config.screnHeight, Phaser.CANVAS, 'content', null);

    this.state.add('Boot', BootState, false);
    this.state.add('Splash', SplashState, false);
    this.state.add('Game', GameState, false);
  }
}

function launchGame () {
  window.game = new Game();
  document.getElementById('startMenuWrapper').hidden = true;
  window.game.state.playerName = playerNameInput.value;
  window.game.state.start('Boot');
}

// Setup start Button
var btn = document.getElementById('startButton');
var playerNameInput = document.getElementById('playerNameInput');

// If launch button is clicked => lauch game
btn.onclick = function () {
  launchGame();
};

playerNameInput.addEventListener('keypress', function (e) {
  var key = e.which || e.keyCode;
  if (key === 13) {
    launchGame();
  }
});
