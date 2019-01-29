import Phaser from 'phaser'

// We only have one phase (game screen)
// If we have more than one phase, we have to generate new State class
export default class extends Phaser.State {
  init () {}

  preload () {
    // load your assets
    this.load.image('bullet', 'assets/images/bullet.png');
    this.load.image('maptile', 'assets/images/maptile.png');
    this.load.image('vietnam', 'assets/images/vietnam.png');
    this.load.image('player', 'assets/images/player.png');
    this.load.image('enemy', 'assets/images/enemy.png');
    this.load.image('wall', 'assets/images/wall.png');
    this.load.image('player_particle', 'assets/images/particle/player_particle.png');
  }

  create () {
    this.state.start('Game');
  }
}
