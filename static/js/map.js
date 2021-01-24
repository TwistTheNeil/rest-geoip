function renderMap(mymap, coordinates, mapboxAccessToken) {
  var mapboxTheme = currentTheme() === "dark" ? "mapbox/dark-v10" : "mapbox/streets-v11";

  L.marker(coordinates).addTo(mymap);

  L.tileLayer(`https://api.mapbox.com/styles/v1/{id}/tiles/{z}/{x}/{y}?access_token=${mapboxAccessToken}`, {
    maxZoom: 18,
    attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
      '<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
      'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
    id: mapboxTheme, // Get these style IDs here https://docs.mapbox.com/api/maps/styles/#mapbox-styles
    tileSize: 512,
    zoomOffset: -1
  }).addTo(mymap);
}
