'use strict';

/**
 * @ngdoc function
 * @name HackerNewsTabApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the HackerNewsTabApp
 */
angular.module('HackerNewsTabApp')
  .controller('MainCtrl', function ($scope, HackerNews) {
    $scope.topStories =
      HackerNews.fetchHomepage();
  });
