const CONTENT_MESSAGE="content_message"
const SYSTEM_MESSAGE="system_message"
const REGISTRATION_MESSAGE="registration_message"

//AUTHENTICATION
const BASIC = "BASIC"


const REGISTRATION="REGISTRATION"
const MESSAGE="MESSAGE"
const NOTIFICATION="NOTIFICATION"
const SUBSCRIPTION="SUBSCRIPTION"
const UNSUBSCRIPTION="UNSUBSCRIPTION"
const REGISTER="REGISTRATION"
const ERROR="ERROR"
const HEADERS = {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
}

class RTMessage{
    constructor(type,channels,client,content,content_type,comunication){
        this.type=type
        this.channels=channels
        this.comunication=comunication
        this.client=client
        this.content=content
        this.content_type=content_type
    }

    toJson(){
        return{
            type:this.type,
            comunication:this.comunication,
            client:null,
            content:this.content,
            content_type:this.content_type
        }
    }
    toJsonString(){
        return JSON.stringify(this.toJson())
    }
}

class RTNamespace{
    constructor(type,name){
        this.type=type
        this.name=name
    }
    toJson(){
        return{
            type:this.type,
            name:this.name
        }
    }
    toJsonString(){
        return JSON.stringify(this.toJson())
    }
}

class RTComunication{
    constructor(type,namespaces,client){
        this.type=type
        this.namespaces=namespaces
        this.client=client
    }
    toJson(){
        return{
            type:this.type,
            namespaces:this.namespaces,
            client:this.client
        }
    }
    toJsonString(){
        return JSON.stringify(this.toJson())
    }
}

class RTMessageRegister extends RTMessage{
    constructor(){
        let communication = new RTComunication(REGISTER,null,null)
        super(MESSAGE,null,null,null,CONTENT_MESSAGE,communication)
    }
}
class RTMessageSubscription extends RTMessage{
    constructor(namespace_names){
        let nodespace = new RTNamespace("PUBLIC",namespace_names)
        let nodespace_array = [nodespace]
        let communication = new RTComunication(SUBSCRIPTION,nodespace_array,null)
        super(MESSAGE,null,null,null,CONTENT_MESSAGE,communication)
    }
}

class RTClient {
    constructor(url,port,messageEvent,closeEvent,saveEvent,updateEvent) {
        this.namespaces_subscribed = []
        this.messageEvent = messageEvent
        this.saveEvent=saveEvent
        this.updateEvent=updateEvent
        this.closeEvent = closeEvent
        this.idClient = null
        this.messageList = []
        this.ws =  new WebSocket(`ws://${url}:${port}`)

        this.ws.addEventListener('open',(event)=> {
            let register_message = new RTMessageRegister()
            this.ws.send(register_message.toJsonString())
            //this.dispatchEvent(new Event('change'));
            this.saveEvent()
        });
        this.ws.addEventListener('message',(event) =>{
            let json_message = JSON.parse(event.data)
            this.messageRouter(json_message)
            this.messageEvent(json_message)
        })
        this.ws.addEventListener('close',(event) =>{
            this.close()
        })
    }

    subscription(namespaces){
        let message = new RTMessageSubscription(namespaces)

    }

    message(message,namespaces){
        let body =  JSON.stringify({
            "timestamp":(new Date()).toISOString(),
            "type":MESSAGE,
            "comunication":{
                "type":MESSAGE,
                "namespaces":namespaces
            },
            "content":message
        })
        fetch('http://localhost:8084/api/message',{
            method: 'POST',
            body:body,
            headers:this.headers}).then(function(response) {
            console.log(response)

        })
            .then(function(response) {
                console.log(response)

            })
    }

    close(){

        this.ws.close()
        this.closeEvent()
    }

    messageRouter(message){
        switch(message.type){
            case MESSAGE:
                if (message.type == "content_message"){
                    this.receptionMessage(message)
                }
            case NOTIFICATION:
                if (message.comunication.type == REGISTRATION){
                    this.idClient = message.comunication.client.uuid
                    this.updateEvent()
                }
                else if(message.comunication.type == SUBSCRIPTION){
                    console.log(message.comunication.namespaces)
                    this.namespaces_subscribed = this.namespaces_subscribed.concat(message.comunication.namespaces)
                    this.updateEvent()
                }
                else if(message.comunication.type == UNSUBSCRIPTION){
                    console.log(message.comunication.namespaces)
                    this.namespaces_subscribed.splice(this.namespaces_subscribed.findIndex(e => e === message.comunication.namespaces))
                    this.updateEvent()
                }
        }

    }

    send(message){
        this.ws.send(message)
    }
}


class RTAuthentication{
    constructor(type,token){
        this.type = type
        this.token = token
    }

    getToken(){
        if (this.type === BASIC){
            return `BASIC ${this.token}`
        }else{
            return `${this.token}`
        }
    }

}


class RTApiClient{
    constructor(url,port,authentication) {
        this.authentication = authentication
        this.idClient = null
        this.messageList = []
        this.host =  `http://${url}:${port}`
        this.headers = HEADERS
        if(this.authentication != null){
            this.headers["AUTHORIZATION"]=`${this.authentication}`

        }
    }
    subscription(namespaces,idClient,callback){
        let resource= "/api/subscribe"
        let body =  JSON.stringify({
            "timestamp":(new Date()).toISOString(),
            "type":MESSAGE,
            "comunication":{
                "type":SUBSCRIPTION,
                "namespaces":namespaces,
                "client":{
                    "uuid":idClient
                }
            },
        })
        fetch(`${this.host}${resource}`,{
            method: 'POST',
            body:body,
            headers:this.headers}
        ).then(function(response) { return response.json(); })
            .then(function(data) {
                callback(data)
            })
    }

    unsubscription(namespaces,idClient,callback){
        let resource= "/api/unsubscribe"
        let body =  JSON.stringify({
            "timestamp":(new Date()).toISOString(),
            "type":MESSAGE,
            "comunication":{
                "type":UNSUBSCRIPTION,
                "namespaces":namespaces,
                "client":{
                    "uuid":idClient
                }
            },
        })
        fetch(`${this.host}${resource}`,{
            method: 'POST',
            body:body,
            headers:this.headers}
        ).then(function(response) { return response.json(); })
            .then(function(data) {
                callback(data)
            })
    }

    list_namespace(callback){
        let resource="/api/namespaces"
        fetch(`${this.host}${resource}`,{
            method: 'GET',
            headers:this.headers}
        ).then(function(response) { return response.json(); })
            .then(function(data) {
                callback(data)
            })
    }


    list_connections(callback){
        let resource="/api/clients"
        return fetch(`${this.host}${resource}`,{
            method: 'GET',
            headers:this.headers}
        ).then(function(response) { return response.json(); })
        .then(function(data) {
                callback(data)
            })
    }
    close(id,callback){
        let resource=`/api/clients/${id}`
        return fetch(`${this.host}${resource}`,{
            method: 'DELETE',
            headers:this.headers}
        ).then(function(response) { return response.json(); })
            .then(function(data) {
                callback(data)
            })
    }

    message(content,namespaces,callback){
        let resource= "/api/message"
        let body =  JSON.stringify({
            "timestamp":(new Date()).toISOString(),
            "type":MESSAGE,
            "comunication":{
                "type":MESSAGE,
                "namespaces":namespaces
            },
            "content":content
        })
        fetch(`${this.host}${resource}`,{
            method: 'POST',
            body:body,
            headers:this.headers})
            .then(function(response) { return response.json(); })
            .then(function(data) {
                callback(data)
            })
    }
}

export {
    RTClient,
    RTApiClient,
    RTAuthentication,
    BASIC
}
