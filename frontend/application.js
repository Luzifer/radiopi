// Generated by CoffeeScript 1.9.3
(function() {
  var loadFavorites, loadVersion, playButtonClick, searchInputType, stopButtonClick, submitPlay, urlListClick;

  $(function() {
    loadVersion();
    $('#play').bind('click', playButtonClick);
    $('#stop').bind('click', stopButtonClick);
    $('#searchInput').bind('keydown', function() {
      clearTimeout(window.searchInputTimer);
      return window.searchInputTimer = setTimeout(searchInputType, 1000);
    });
    return loadFavorites();
  });

  stopButtonClick = function() {
    submitPlay('off');
    return false;
  };

  playButtonClick = function() {
    var url;
    url = $('#urlInput').val();
    submitPlay(url);
    return false;
  };

  submitPlay = function(url) {
    return $.post('/v1/play', {
      stream: url
    }, function(data) {
      return $('#urlInput').val('');
    });
  };

  searchInputType = function() {
    var search;
    search = $('#searchInput').val();
    if (search.length < 3) {
      return;
    }
    search = encodeURIComponent(search);
    return $.get("/v1/search?search=" + search, function(data) {
      var a, entry, i, len, results;
      $('#stationList').empty();
      results = [];
      for (i = 0, len = data.length; i < len; i++) {
        entry = data[i];
        a = $("<a>" + entry.server_name + "</a>");
        a.attr('href', 'javascript:void(0)');
        a.attr('class', 'list-group-item');
        a.data('url', entry.listen_url);
        a.bind('click', urlListClick);
        results.push(a.appendTo($('#stationList')));
      }
      return results;
    });
  };

  urlListClick = function() {
    var url;
    url = $(this).data('url');
    return submitPlay(url);
  };

  loadFavorites = function() {
    return $.get("/v1/favorites", function(data) {
      var a, entry, i, len, results;
      if (data.length > 0) {
        $('#favoritePanel').removeClass("hidden");
      } else {
        $('#favoritePanel').addClass("hidden");
      }
      $('#favoriteList').empty();
      results = [];
      for (i = 0, len = data.length; i < len; i++) {
        entry = data[i];
        a = $("<a>" + entry.name + "</a>");
        a.attr('href', 'javascript:void(0)');
        a.attr('class', 'list-group-item');
        a.data('url', entry.url);
        a.bind('click', urlListClick);
        results.push(a.appendTo($('#favoriteList')));
      }
      return results;
    });
  };

  loadVersion = function() {
    return $.get("/v1/version", function(data) {
      return $("#version").text(data);
    });
  };

}).call(this);
