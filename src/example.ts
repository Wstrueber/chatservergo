const socket = new WebSocket("ws://localhost:8080");
socket.onopen = () => {
  console.log("Successfully Connected");
  socket.send(JSON.stringify({ action: "REQUEST_VERSION_NUMBER" }));
};
socket.onclose = e => {
  console.log("Socket Closed Connection: ", e);
};
socket.onmessage = message => {
  console.log(JSON.parse(message.data));
};
socket.onerror = e => {
  console.log("Socket Error: ", e);
};

socket.send(
  JSON.stringify({ action: "REQUEST_RESPONSE", message: "Hello Server" })
);
