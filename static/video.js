const canvas = document.createElement("canvas");
const wsavc = new WSAvcPlayer(canvas, "webgl", 1, 35);

wsavc.connect("ws://" + document.location.host + "/video");
document.getElementById('video-container').appendChild(canvas);
setTimeout(function() { wsavc.playStream() }, 3000);
