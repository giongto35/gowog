# GOWOG Frontend

## Phaser + ES6 + Webpack.
![Phaser+ES6+Webpack](https://raw.githubusercontent.com/lean/phaser-es6-webpack/master/assets/images/phaser-es6-webpack.jpg)

# Setup
You'll need to install a few things before you have a working copy of the project.


```git clone https://github.com/lean/phaser-es6-webpack.git```

## 1. Install node.js and npm:

https://nodejs.org/en/

## 2. Install dependencies (optionally you can install [yarn](https://yarnpkg.com/)):

```npm install``` 

or if you chose yarn, just run ```yarn```

## 3. Run the development frontend:

Run:

```npm run dev -- --env.HOST_IP=localhost:8080```

This will run web game client in the browser. It will also start a watch process, so you can change the source and the process will recompile and refresh the browser automatically.
Note: localhost:8080 is the address of webserver host.

To see the game, open your browser and enter http://localhost:3000 into the address bar.

# Credits
The Frontend codebase is based on
https://github.com/RenaudROHLINGER/phaser-es6-webpack
