'use strict';

/**
 * @ngdoc function
 * @name hackerNewsTabApp.controller:MainCtrl
 * @description
 * # MainCtrl
 * Controller of the hackerNewsTabApp
 */
angular.module('hackerNewsTabApp')
  .controller('MainCtrl', function ($scope) {
    $scope.awesomeThings = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ];
  });
