# Golang with Gin & Melody as PingPong

> Measure latency of websocket ping-pong to the server

> This actually not ping-pong mechanism, because ping-pong are handled automatically by client-server frame. This implementation just to take advantage of `HandlePong` on Melody module then measure the timestamps that sended by client to server (round trip).

![image](https://user-images.githubusercontent.com/738088/145684183-9a0f9988-2ac6-4de3-bdf1-6271f7947534.png)
