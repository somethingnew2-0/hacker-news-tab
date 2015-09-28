'use strict';

angular.module('HackerNewsTabApp')
  .filter('image', function() {
    return function(url) {
      if ('undefined' !== typeof url) {
        return '/screenshot?url='+encodeURIComponent(url);
      }
    };
  });
