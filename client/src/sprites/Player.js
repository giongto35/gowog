import Phaser from 'phaser';
import Shoot from './Shoot';
import config from '../config';

export default class extends Phaser.Sprite {
  constructor ({ game, layer, id, name, x, y, asset }) {
    super(game, x, y);

    this.id = id;
    this.x = x;
    this.y = y;
    this.size = config.playerSize;
    this.name = name;
    this.nextReload = 0;
    this.score = 0;

    // Setup circle graphic for player
    var graphics = new Phaser.Graphics(game, 0, 0);
    graphics.lineStyle(2, 0x000000);
    graphics.beginFill(0x0000FF, 1);
    graphics.drawCircle(0, 0, config.playerSize);
    this.addChild(graphics);

    this.shootManager = new Shoot({
      game: game,
      layer: layer,
      x: x,
      y: y,
      key: 'bullet',
      frame: 0,
      playerID: id,
      asset: null
    });
    this.anchor.setTo(0.5, 0.5);

    game.add.existing(this);

    // Healthbar
    this.health = 100.0;
    this.healthbar = game.add.graphics(0, 0);
    this.healthbar.anchor.setTo(0.5, 0.5);
    this.healthbar.beginFill(0x00ff00);
    this.healthbar.drawRect(-50, -40, 100, 15);
    this.addChild(this.healthbar);

    // Name
    this.nametag = game.add.text(0, 0, this.name);
    this.nametag.stroke = '#000000';
    this.nametag.strokeThickness = 1;
    this.nametag.anchor.setTo(0.5, 2.3);
    this.nametag.fill = '#ffffff';
    this.nametag.setShadow(1, 2, '#333333', 2, false, true);
    this.addChild(this.nametag);
    this.nametag.bringToTop();

    // this.turret.scale.setTo(game.scaleRatio, game.scaleRatio);
    this.game.physics.enable(this, Phaser.Physics.ARCADE);

    layer.add(this);
  }

  update () {
    this.healthbar.width = this.health / 100 * 150;
  }

  fire (x, y, dx, dy) {
    //  Grab the first bullet we can from the pool
    var shoot = this.shootManager.getFirstExists(false);

    if (shoot) {
      // Fire it
      shoot.reset(x, y);
      shoot.body.velocity.y = dy * 1000;
      shoot.body.velocity.x = dx * 1000;
      var deg = Math.atan2(dx, dy);
      shoot.angle = deg * -180 / Math.PI + 90;
      console.log(deg, shoot.angle);
    }
  }

  move (moveRight, moveDown) {
    // Params with both moveRight and moveDown, so we can know the direction player is facing
    // For this game, we don't need facing direction
    // TODO: Move animation here
  }
}
