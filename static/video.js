let frames = [];
let framesRendered = 0;

const videoSocket = new WebSocket(`ws://${location.host}/video`);
videoSocket.binaryType = 'arraybuffer';

const videoContainer = document.getElementById('video-container');
const player = new Player({ 
  useWorker: true, 
  workerFile: './broadway/Decoder.js',
  webgl: true
});
videoContainer.appendChild(player.canvas);

videoSocket.addEventListener('open', function (e) {
  console.log('Video socket connected');
  videoSocket.send(JSON.stringify({ action: "start" }));
  render();
});

videoSocket.addEventListener('close', function (e) {
  console.log('Video socket closed');
});

videoSocket.addEventListener('message', function (e) {
  if (typeof e.data !== 'string') frames.push(e);
});

function render() {
  if (framesRendered > 10 && frames.length > 10){
    console.log('dropping frames');
    frames = frames.slice(10);
  }

  if (frames.length > 0) {
    const frame = frames.shift();
    player.decode(new Uint8Array(frame.data, 0, frame.data.size));
  }

  framesRendered++;
  requestAnimationFrame(render);
}
