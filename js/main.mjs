import net from "net";

const conns = [];

net.createServer((socket) => {
    conns.push(socket);

    socket.on("data", (data) => {
        for (const conn of conns) {
            conn.write(data);
        }
    });

    socket.on("close", () => {
        conns.splice(conns.indexOf(socket), 1);
    });
}).listen(1337);
