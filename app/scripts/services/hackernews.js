'use strict';

/**
 * @ngdoc service
 * @name HackerNewsTabApp.HackerNews
 * @description
 * # HackerNews
 * Service in the hackerNewsTabApp.
 */
angular.module('HackerNewsTabApp')
  .service('HackerNews', ['Ref', '$firebaseArray', '$firebaseObject', function (Ref, $firebaseArray, $firebaseObject) {
    return { 
      fetchHomepage: function() {
        return $firebaseArray(Ref.child('topstories/').orderByKey().limitToFirst(25));
      },
      fetchItem: function(itemId) {
        return $firebaseObject(Ref.child('item/' + itemId));
      },
      fetchUser: function(userId) {
        return $firebaseObject(Ref.child('user/' + userId));
      }
    };
  }]);
