const player = new JSMpeg.Player(`ws://${location.host}/video`, {
  canvas: document.getElementById('video-container'),
  audio: false,
  autoplay: true,
});
