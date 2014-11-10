/*global define */

define(['Phaser', 'StartState'], function (Phaser, StartState) {
    'use strict';

    function Bumboo(id) {
        Phaser.Game.call(this, 800, 600, Phaser.AUTO, id);

        var clearWorld = true,
            clearCache = true,
            autoStart = true;

        this.boot();
        this.stage.setBackgroundColor('#3498DB');

        this.state.add('Start', StartState, !autoStart);
        this.state.start('Start', clearWorld, !clearCache);
    }

    Bumboo.prototype = Object.create(Phaser.Game.prototype);

    return Bumboo;
});