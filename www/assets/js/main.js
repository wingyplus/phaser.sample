/*global require */

require.config({
    baseUrl: 'assets/js',
    paths: {
        Phaser: 'libs/phaser/build/phaser',
        domReady: 'libs/requirejs-domready/domReady',
        Player: 'Player',
        Bumboo: 'Bumboo'
    }
});

require(['domReady', 'Bumboo'], function (domReady, Bumboo) {
    'use strict';

    domReady(function () {
        new Bumboo();
    });
});