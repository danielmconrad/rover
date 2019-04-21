// const videoSocket = new WebSocket(`ws://${location.host}/video`);
// const player = new Player({
//   size: {
//     width: 320,
//     height: 240,
//   },
// });

// document.getElementById('video-container').appendChild(player.canvas);

// videoSocket.addEventListener('open', function (e) {
//   console.log('Socket connected');
// });

// videoSocket.addEventListener('close', function (e) {
//   console.log('Socket closed');
// });

// videoSocket.addEventListener('message', function (e) {
//   player.decode(new Uint8Array(e.data));
// });

const canvas = document.createElement("canvas");
document.getElementById('video-container').appendChild(canvas);

const wsavc = new WSAvcPlayer(canvas, "webgl", 1, 35);
wsavc.connect("ws://" + document.location.host + "/video");
wsavc.playStream();