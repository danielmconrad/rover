const videoSocket = new WebSocket(`ws://${location.host}/video`);
const videoContainer = document.getElementById('video-container');

var decoder = new Decoder();
document.getElementById('video-container').appendChild(decoder.canvas);

videoSocket.addEventListener('open', function (e) {
  console.log('Video socket connected');
});

videoSocket.addEventListener('close', function (e) {
  console.log('Video socket closed');
});

videoSocket.addEventListener('message', function (e) {
  console.log('Video message from server ', e.data);
  decoder.decode(new Uint8Array(e.data));
});
