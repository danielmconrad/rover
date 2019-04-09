$(document).foundation();

const socket = new WebSocket(`ws://${location.host}/messages`);

socket.addEventListener('open', function (event) {
  socket.send(JSON.stringify({ 
    event: 'first' 
  }));
  socket.send(JSON.stringify({ 
    event: 'second' 
  }));
  socket.send(JSON.stringify({ 
    event: 'third' 
  }));
});

socket.addEventListener('close', function (event) {
  console.log('Socket closed');
});

socket.addEventListener('message', function (event) {
  const message = JSON.parse(event.data);
  console.log('Message from server ', message);
});