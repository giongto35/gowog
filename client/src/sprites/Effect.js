import Phaser from 'phaser';

export default function explode (game, layer, x, y) {
  var emitter = game.add.emitter(x, y, 200);
  emitter.makeParticles('explode_particle');
  emitter.gravity = 0;
  emitter.setScale(0.2, 0.1, 0.2, 0.1, 500, Phaser.Easing.Linear.None);
  emitter.setAlpha(0.9, 0.6, 500, Phaser.Easing.Linear.None, false);
  emitter.setXSpeed(-1000, 1000);
  emitter.setYSpeed(-1000, 1000);
  emitter.explode(300, 200);
  layer.add(emitter);
}
