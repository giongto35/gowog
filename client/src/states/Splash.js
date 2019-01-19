import Phaser from 'phaser'

// We only have one phase (game screen)
// If we have more than one phase, we have to generate new State class
export default class extends Phaser.State {
  init () {}

  preload () {
    // load your assets
    this.load.image('bullet', 'assets/images/bullet.png');
    this.load.image('maptile', 'assets/images/maptile.png');
  }

  create () {
    this.state.start('Game');
  }
}
