const socket = new WebSocket("ws://localhost:8080/ws");
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

socket.send(
  JSON.stringify({
    action: "REQUEST_LOGIN",
    client: {
      clientId: "3f43c165-e9f6-47ca-83b8-9daf00bafc57",
      userName: "billy"
    }
  })
);
