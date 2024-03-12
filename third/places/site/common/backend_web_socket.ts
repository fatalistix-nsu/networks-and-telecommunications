// for tls insert your own dns
// const domainName = "your-dns"
const domainName = "localhost"
const webSocketListenURI = "/ws"
// only for TLS
// export const ws = new WebSocket("wss://" + domainName + webSocketListenURI)
export const ws = new WebSocket("ws://" + domainName + webSocketListenURI)

export type ServerResponse = {
    code: number,
    info: any
}
