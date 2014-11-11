/*global define */

define(['Phaser', 'Player'], function (Phaser, Player) {
    'use strict';

    function StartState() {
        this.player = null;
    }

    StartState.prototype = Object.create(Phaser.State.prototype);

    StartState.prototype.preload = function () {
        this.load.image('player', 'assets/player.png');
    };

    StartState.prototype.create = function () {
        this.player = new Player(this, this.world.centerX, this.world.centerY, 'player');
        this.physics.enable(this.player, Phaser.ARCADE);
        this.add.existing(this.player);
    };

    StartState.prototype.update = function () {
        if (this.game.input.keyboard.isDown(Phaser.Keyboard.LEFT)) {
            this.player.moveLeft();
        } else if (this.game.input.keyboard.isDown(Phaser.Keyboard.RIGHT)) {
            this.player.moveRight();
        } else {
            this.player.dontMove();
        }
    };

    return StartState;
});
