class WebsocketStore{
	constructor(){
        this.websocket_list = new Array()
    }
    addConnection(connection){
    	this.websocket_list.push(connection)
    }
}

export default WebsocketStore