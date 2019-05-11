const videoURL = `ws://${location.host}/video`;
const videoSocket = new WebSocket(videoURL);
const videoContainer = document.getElementById('video-container');
const videoCanvas = document.getElementById('video-canvas');

var decoder = new Decoder();
var running = false;
var frames = [];

videoSocket.addEventListener('open', function (e) {
  console.log('Video socket connected');
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
    
    videoSocket.send(JSON.stringify({ action: "start" }));
    render();
  }
});

function render() {
  if (frames.length) { 
    decoder.decode(new Uint8Array(frames.shift()));
  }
  
  requestAnimationFrame(render);
}
