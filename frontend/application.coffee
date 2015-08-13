$ ->
  loadVersion()
  $('#play').bind 'click', playButtonClick
  $('#stop').bind 'click', stopButtonClick
  $('#searchInput').bind 'keydown', () ->
    clearTimeout(window.searchInputTimer)
    window.searchInputTimer = setTimeout searchInputType, 1000
  loadFavorites()

stopButtonClick = () ->
  submitPlay 'off'
  false

playButtonClick = () ->
  url = $('#urlInput').val()
  submitPlay url
  false

submitPlay = (url) ->
  $.post '/v1/play',
    stream: url
  , (data) ->
    $('#urlInput').val ''

searchInputType = () ->
  search =   $('#searchInput').val()
  if search.length < 3
    return

  search = encodeURIComponent(search)
  $.get "/v1/search?search=#{search}", (data) ->
    $('#stationList').empty()
    for entry in data
      a = $("<a>#{entry.server_name}</a>")
      a.attr 'href', 'javascript:void(0)'
      a.attr 'class', 'list-group-item'
      a.data 'url', entry.listen_url
      a.bind 'click', urlListClick
      a.appendTo $('#stationList')

urlListClick = () ->
  url = $(this).data 'url'
  submitPlay url

loadFavorites = () ->
  $.get "/v1/favorites", (data) ->
    if data.length > 0
      $('#favoritePanel').removeClass("hidden")
    else
      $('#favoritePanel').addClass("hidden")
    $('#favoriteList').empty()
    for entry in data
      a = $("<a>#{entry.name}</a>")
      a.attr 'href', 'javascript:void(0)'
      a.attr 'class', 'list-group-item'
      a.data 'url', entry.url
      a.bind 'click', urlListClick
      a.appendTo $('#favoriteList')

loadVersion = () ->
  $.get "/v1/version", (data) ->
    $("#version").text(data)
