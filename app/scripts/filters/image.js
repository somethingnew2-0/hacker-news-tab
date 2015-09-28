'use strict';

angular.module('HackerNewsTabApp')
  .filter('image', function() {
    return function(url) {
      if ('undefined' !== typeof url) {
        return 'http://localhost:8000/screenshot?url='+encodeURIComponent(url);
      }
    };
  });
