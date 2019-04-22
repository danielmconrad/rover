const canvas = document.createElement("canvas");
document.getElementById('video-container').appendChild(canvas);

const wsavc = new WSAvcPlayer(canvas, "webgl", 1, 35);
wsavc.connect("ws://" + document.location.host + "/video");

setTimeou(function() {
    wsavc.playStream();
}, 2000);