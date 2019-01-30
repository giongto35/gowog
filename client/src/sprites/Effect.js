import Phaser from 'phaser';

export function explode (game, layer, x, y) {
  var emitter = game.add.emitter(x, y, 200);
  emitter.makeParticles('explode_particle');
  emitter.gravity = 0;
  emitter.setScale(0.2, 0.1, 0.2, 0.1, 500, Phaser.Easing.Linear.None);
  emitter.setAlpha(0.9, 0.6, 500, Phaser.Easing.Linear.None, false);
  emitter.setXSpeed(-1000, 1000);
  emitter.setYSpeed(-1000, 1000);
  emitter.explode(300, 200);
  layer.add(emitter);

  game.time.events.add(200, function () { emitter.destroy(); }, this);
}

export function explode_bullet (game, layer, x, y) {
  var emitter = game.add.emitter(x, y, 40);
  emitter.makeParticles('explode_bullet');
  emitter.gravity = 0;
  emitter.setScale(0.25, 0.15, 0.25, 0.15, 150, Phaser.Easing.Linear.None);
  emitter.setAlpha(0.8, 0.6, 150, Phaser.Easing.Linear.None, false);
  emitter.setXSpeed(-1000, 1000);
  emitter.setYSpeed(-1000, 1000);
  emitter.explode(150, 40);
  layer.add(emitter);

  game.time.events.add(200, function () { emitter.destroy(); }, this);
}
