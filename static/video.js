const videoSocket = new WebSocket(`ws://${location.host}/video`);
const videoContainer = document.getElementById('video-container');

var decoder = new Decoder();
var running = false;
var frames = [];

videoSocket.addEventListener('open', function (e) {
  console.log('Video socket connected');
  videoSocket.send(JSON.stringify({ action: "start" }));
  render();
});

videoSocket.addEventListener('close', function (e) {
  console.log('Video socket closed');
});

videoSocket.addEventListener('message', function (e) {
  if (typeof e.data !== 'string') 
    return frames.push(e.data);

  const message = JSON.parse(e.data);

  if (message.action === 'init') {
    const canvasWrapper = new YUVCanvas(message);
    decoder.onPictureDecoded = canvasWrapper.decode;
    videoContainer.appendChild(canvasWrapper.canvasElement);
  }
});

function render() {
  if (!running) return;
  if (frames.length >= 10) {
    console.log('dropping')
    return frames = [];
  }
  return decoder.decode(new Uint8Array(frames.pop()));
  requestAnimationFrame(render);
}
