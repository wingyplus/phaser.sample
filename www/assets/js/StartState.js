/*global define */

define(['Phaser', 'Player'], function (Phaser, Player) {
    'use strict';

    function StartState() {
        this.player = null;
        /*
         * @param {List<Player>}
         */
        this.otherPlayers = [];
        this.ws = null;
    }

    StartState.prototype = Object.create(Phaser.State.prototype);

    StartState.prototype.preload = function () {
        this.load.image('player', 'assets/player.png');
    };

    StartState.prototype.create = function () {
        // make a connection between client and server.
        var ws = new WebSocket('ws://localhost:8100/join');
        ws.onopen = function onOpen(e) {
            ws.send(JSON.stringify({
                type:'create'
            }));
        };
        var onMessageFunc = (function (state) {
            return function onMessage(e) {
                var action = JSON.parse(e.data);
                if (action.type === 'create') {
                    console.log(action);
                    state.player = new Player(state, action.id, state.world.centerX, state.world.centerY, 'player');
                    state.physics.enable(state.player, Phaser.ARCADE);
                    state.add.existing(state.player);
                } else if (action.type === 'otherCreate') {
                    var p = new Player(state, action.id, action.x, action.y, 'player')
                    state.otherPlayers.push(p);
                    state.physics.enable(p, Phaser.ARCADE);
                    state.add.existing(p);
                    console.log(p);
                }
            };
        })(this);
        ws.onmessage = onMessageFunc;
        this.ws = ws;
    };

    StartState.prototype.update = function () {
        if (!this.player) {
            // do nothing
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.LEFT)) {
            this.player.moveLeft();
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.RIGHT)) {
            this.player.moveRight();
        } else {
            this.player.dontMove();
        }
    };

    return StartState;
});
