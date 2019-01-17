import Phaser from 'phaser';
import config from '../config';

export default class extends Phaser.Sprite {
  constructor ({game, x, y, width, height}) {
    super(game, 0, 0);

    this.topList = [];
    this.numTop = config.numTop;

    var rect = new Phaser.Graphics(game, 0, 0);
    rect.lineStyle(2, 0x000000);
    rect.beginFill(0x000000, 1);
    rect.drawRect(config.screenWidth - config.leaderboardSize, 0, config.leaderboardSize, this.numTop * 25 + 10);
    rect.alpha = 0.3;
    this.addChild(rect);

    for (var i = 0; i < this.numTop; i++) {
      var text = game.add.text(config.screenWidth - config.leaderboardSize + 10, i * 25, '');
      text.bringToTop();
      this.topList.push(text);
      this.addChild(text);
    }

    game.add.existing(this);
  }

  updateLeaderboard (topPlayers) {
    for (var i = 0; i < Math.min(topPlayers.length, this.numTop); i++) {
      this.topList[i].text = topPlayers[i].name + '  ' + topPlayers[i].score;
    }
  }
}
