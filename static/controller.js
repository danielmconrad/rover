const controllerSocket = new WebSocket(`ws://${location.host}/controller`);
const infoContainer = document.getElementById("info-container");
const controllerMap = new Map();

const SOCKET = {
  CONNECTING: 0,
  OPEN: 1,
  CLOSING: 2,
  CLOSED: 3
};

controllerSocket.addEventListener("open", function(e) {
  console.log("Controller socket connected");
  sendControllerState();
});

controllerSocket.addEventListener("close", function(e) {
  console.log("Controller socket closed");
});

controllerSocket.addEventListener("message", function(e) {
  console.log("Controller message from server ", JSON.parse(e.data));
});

window.addEventListener("gamepadconnected", function(e) {
  console.log("Gamepad connected:", e.gamepad);
  controllerMap.set(e.gamepad, e.gamepad);

  e.gamepad.vibrationActuator.playEffect("dual-rumble", {
    startDelay: 0,
    duration: 200,
    weakMagnitude: 1.0,
    strongMagnitude: 1.0
  });
});

window.addEventListener("gamepaddisconnected", function(e) {
  console.log("Gamepad disconnected:", e.gamepad);
  controllerMap.delete(e.gamepad);
});

function sendControllerState() {
  if (controllerSocket.readyState !== SOCKET.OPEN) return;

  for (let controller of controllerMap.values()) {
    const gamepad = navigator.getGamepads()[controller.index];

    infoContainer.innerHTML = `
      <span>
        ${gamepad.axes[1].toFixed(2)}
        ${gamepad.axes[3].toFixed(2)}
      </span>
    `;

    controllerSocket.send(
      JSON.stringify({
        left: gamepad.axes[1],
        right: gamepad.axes[3]
      })
    );
  }

  requestAnimationFrame(sendControllerState);
}
