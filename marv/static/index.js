$(document).foundation();

const controllerSocket = new WebSocket(`ws://${location.host}/controller`);

var firstMessageDate;

controllerSocket.addEventListener('open', function (event) {
  controllerSocket.send(JSON.stringify({ 
    event: 'first' 
  }));
  firstMessageDate = new Date();
});

controllerSocket.addEventListener('close', function (event) {
  console.log('Socket closed');
});

controllerSocket.addEventListener('message', function (event) {
  const now = new Date();
  const message = JSON.parse(event.data);
  console.log('Message from server ', message, now - firstMessageDate);
});