const controllerSocket = new WebSocket(`ws://${location.host}/controller`);
const controllerContainer = document.getElementById('controller-container');

const SOCKET = {
  CONNECTING: 0,
  OPEN: 1,
  CLOSING: 2,
  CLOSED: 3,
};

controllerSocket.addEventListener('open', function (e) {
  console.log('Controller socket connected');
  sendControllerState();
});

controllerSocket.addEventListener('close', function (e) {
  console.log('Controller socket closed');
});

controllerSocket.addEventListener('message', function (e) {
  console.log('Controller message from server ', JSON.parse(e.data));
});

function sendControllerState() {
  if (controllerSocket.readyState !== SOCKET.OPEN) return;

  Object.values(navigator.getGamepads()).forEach(function(gamepad) {
    if (!gamepad) return;

    controllerContainer.innerHTML = `${gamepad.axes[1]} ${gamepad.axes[3]}`;

    controllerSocket.send(JSON.stringify({ 
      left: gamepad.axes[1],
      right: gamepad.axes[3],
    }));
  });

  requestAnimationFrame(sendControllerState);
}
