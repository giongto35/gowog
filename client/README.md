# GOWOG Frontend

GOWOG front end uses *Phaser* + *ES6* + *Webpack*.
![Phaser+ES6+Webpack](https://raw.githubusercontent.com/lean/phaser-es6-webpack/master/assets/images/phaser-es6-webpack.jpg)

# Background

GOWOG client is written in Phaser JS. Phaser is a the leading game engine on web platform. In GOWOG, phaser is mainly used for render purpose because game logic should be done on backend.

# Installation
[**In Main Page**](..)

This will run web game client in the browser. It will also start a watch process, so you can change the source and the process will recompile and refresh the browser automatically.
Note: localhost:8080 is the address of webserver host.
  * npm run dev -- --env.HOST_IP=localhost:8080

To see the game, open your browser and enter http://localhost:3000 into the address bar.

# Development

## Code structure
```
├── client
│   ├── index.html
│   ├── src
│   │   ├── config.js: javascript config
│   │   ├── index.html
│   │   ├── main.js
│   │   ├── sprites
│   │   │   ├── Leaderboard.js: Leaderboard object
│   │   │   ├── Map.js: Map object
│   │   │   ├── Player.js: Player object
│   │   │   └── Shoot.js: Shoot object
│   │   ├── states
│   │   │   ├── Boot.js Boot screen
│   │   │   ├── const.js
│   │   │   ├── Game.js: Game master
│   │   │   ├── message_pb.js: Protobuf Message
│   │   │   ├── Splash.js
│   │   │   └── utils.js
│   │   └── utils.js
```

Each object is a class inherited from Sprite.
Player contain shootManager, each time player fires, a new bullet is generated from shoot manager .

# Credits

The Frontend codebase is based on
https://github.com/RenaudROHLINGER/phaser-es6-webpack
