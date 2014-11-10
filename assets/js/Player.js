/*global define */

define(['Phaser'], function (Phaser) {
    'use strict';

    function Player(game, x, y, name) {
        Phaser.Sprite.call(this, game, x, y, name, 0);
        game.physics.enable(this, Phaser.ARCADE);
        game.add.existing(this);
    }

    Player.prototype = Object.create(Phaser.Sprite.prototype);

    Player.prototype.moveLeft = function moveLeft() {
        this.body.velocity.x = -200;
    };

    Player.prototype.moveRight = function moveRight() {
        this.body.velocity.x = 200;
    };

    Player.prototype.dontMove = function dontMove() {
        this.body.velocity.x = 0;
    };

    return Player;
});