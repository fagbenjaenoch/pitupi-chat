import net from "net";

function handleNewConnection(socket) {
	socket.on("connect", () => {
		console.log("info: new socket just connected.");
	});

	socket.on("data", (data) => {
		console.log(`info: data from peer: ${data}\n`);
	});

	socket.on("close", () => {
		console.log("info: connection has been half closed");
	});
}

const server = net.createServer((socket) => {
	handleNewConnection(socket);
});

server.listen("8080", () => {
	console.log("info: server listening on port 8080");
});
