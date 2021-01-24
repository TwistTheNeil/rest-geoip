function changeTheme() {
  switch (currentTheme()) {
    case "light":
      document.cookie = "theme=dark";
      $('#dark-theme-stylesheet').removeAttr('disabled');
      break;
    case "dark":
      document.cookie = "theme=light";
      $('#dark-theme-stylesheet').attr('disabled', 'true');
      break;
  }
}

function currentTheme() {
  return splitCookieString().get('theme') || "light"
}

function splitCookieString() {
  var s = new Map();
  document.cookie.split(';').forEach(function (e) {
    var kv = e.trim().split('=');
    s.set(kv[0], kv[1]);
  });
  return s;
}
