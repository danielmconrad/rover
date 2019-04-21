const videoSocket = new WebSocket(`ws://${location.host}/video`);
const player = new Player();

videoSocket.addEventListener('open', function (e) {
  console.log('Socket connected');
});

videoSocket.addEventListener('close', function (e) {
  console.log('Socket closed');
});

videoSocket.addEventListener('message', function (e) {
  player.decode(e.data);
});