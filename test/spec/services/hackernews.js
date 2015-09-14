'use strict';

describe('Service: HackerNews', function () {

  // load the service's module
  beforeEach(module('hackerNewsTabApp'));

  // instantiate service
  var HackerNews;
  beforeEach(inject(function (_HackerNews_) {
    HackerNews = _HackerNews_;
  }));

  it('should do something', function () {
    expect(!!HackerNews).toBe(true);
  });

});
