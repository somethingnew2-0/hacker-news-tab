'use strict';

angular.module('HackerNewsTabApp')
  .filter('domain', function() {
    return function(url) {
      if ('undefined' !== typeof url) {
        return new URL(url).hostname;
      }
    };
  });
