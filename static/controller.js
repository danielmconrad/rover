const SOCKET = {
  CONNECTING: 0,
  OPEN: 1,
  CLOSING: 2,
  CLOSED: 3,
};

const haveEvents = 'ongamepadconnected' in window;
const controllers = {};
const controllerSocket = new WebSocket(`ws://${location.host}/controller`);
const controllerContainer = document.getElementById('controller-container');

window.addEventListener("gamepadconnected", function(e) {
  controllers[e.gamepad.index] = e.gamepad;
});

window.addEventListener("gamepaddisconnected", function(e) {
  delete controllers[e.gamepad.index];
});

controllerSocket.addEventListener('open', function (e) {
  console.log('Socket connected');
  sendControllerState();
});

controllerSocket.addEventListener('close', function (e) {
  console.log('Socket closed');
});

controllerSocket.addEventListener('message', function (e) {
  // console.log('Message from server ', JSON.parse(e.data));
});

function sendControllerState() {
  if (controllerSocket.readyState !== SOCKET.OPEN) {
    console.log('what happened?')
    return
  }

  Object.values(navigator.getGamepads()).forEach(function(gamepad) {
    if (!gamepad) return;

    controllerContainer.innerHTML = `${gamepad.axes[1]} ${gamepad.axes[3]}`;

    controllerSocket.send(JSON.stringify({ 
      event: 'controller_state',
      data: JSON.stringify({
        left: gamepad.axes[1],
        right: gamepad.axes[3],
      })
    }));
  })
  requestAnimationFrame(sendControllerState)
}