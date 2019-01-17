import Phaser from 'phaser'

// We use phaser group instead of spirite list because phaserGroup support
// reusing. We create a pool of bullets and it will keep using from that pool
export default class extends Phaser.Group {
  constructor ({ game, layer, x, y, key, frame, playerID, asset }) {
    super(game, x, y, key, frame);
    this.playerID = playerID;
    this.enableBody = true;
    this.physicsBodyType = Phaser.Physics.ARCADE;
    // TODO: Set unlimited
    this.createMultiple(30, 'bullet');
    this.setAll('anchor.x', 0.5);
    this.setAll('anchor.y', 0.5);
    this.setAll('outOfBoundsKill', true);
    this.setAll('checkWorldBounds', true);
    game.add.existing(this);
    layer.add(this);
    // this.scale.setTo(game.scaleRatio, game.scaleRatio);
  }

  update () {
    this.forEachAlive(this.updateShoot, this);
  }

  updateShoot (shoot) {
  }

  fire (x, y, dx, dy) {
    //  Grab the first bullet we can from the pool
    var shoot = this.getFirstExists(false);

    if (shoot) {
      //  And fire it
      shoot.reset(x, y);
      shoot.body.velocity.y = dy * 100;
      shoot.body.velocity.x = dx * 100;
    }
  }
}
