'use strict';

/**
 * @ngdoc directive
 * @name HackerNewsTabApp.directive:Story
 * @description
 * # Story
 */
angular.module('HackerNewsTabApp')
  .directive('story', function (HackerNews) {
    return {
      templateUrl: 'views/story.html',
      restrict: 'E',
      link: function link(scope, element, attrs) {
        scope.data = HackerNews.fetchItem(scope.story.$value);
      }
    };
  });
