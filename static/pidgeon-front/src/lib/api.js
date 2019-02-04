import axios from 'axios'

const HOST = "http://192.168.0.136:8084"
const SYSTEM_MESSAGE="system_message"
const CONTENT_MESSAGE="content_message"

const MESSAGE="Message"
const SUBSCRIPTION="Subscription"
const UNSUBSCRIPTION="Unsubscription"
const REGISTER="Register"
const ERROR="ERROR"

class Message{
    constructor(message_type,channels,client,content,content_type){
        this.message_type=message_type
        this.channels=channels
        this.client=client
        this.content=content
        this.content_type=content_type
    }

    toJson(){
        return{
            message_type:this.message_type,
            channels:this.channels,
            client:null,
            content:this.content,
            content_type:this.content_type
        }
    }
    toJsonString(){
        return JSON.stringify(this.toJson())
    }
}



class MessageContent extends Message{
    constructor(channels,content){
        super(MESSAGE,channels,null,content,CONTENT_MESSAGE)
    }
}

class Api{
  constructor() {
  		this.host = HOST
  }
}

class ApiChannel extends Api{
	list(callback){
		axios.get(`${HOST}/api/namespaces`)
	      .then((response) => {
	      	callback(response)
	      })
	      .catch(function (error) {
	    	console.log(error);
	  	  })
	  	  .then(function () {
	  	   });
	}
	retrieve(id,callback){
		axios.get(`${HOST}/api/namespaces/${id}`)
	      .then((response) => {
	      	callback(response)
	      })
	      .catch(function (error) {
	    	console.log(error);
	  	  })
	  	  .then(function () {
	  	   });
	}	
	getClients(id,callback){
		axios.get(`${HOST}/api/namespaces/${id}/clients`)
	      .then((response) => {
	      	callback(response)
	      })
	      .catch(function (error) {
	    	console.log(error);
	  	  })
	  	  .then(function () {
	  	   });
	}


}

class ApiMessage extends Api{

	send(channels,message,callback){
		let message_content = new MessageContent(channels,message)
		axios.post(`${HOST}/api/message`,message_content.toJsonString())
	      .then((response) => {
	      	callback(response)
	      })
	      .catch(function (error) {
	    	console.log(error);
	  	  })
	  	  .then(function () {
	  	   });

	}
}

class ApiClient extends Api{
	list(callback){
		axios.get(`${HOST}/api/clients`)
	      .then((response) => {
	      	callback(response)
	      })
	      .catch(function (error) {
	    	console.log(error);
	  	  })
	  	  .then(function () {
	  	   });
	}
	retrieve(id,callback){
		axios.get(`${HOST}/api/clients/${id}`)
	      .then((response) => {
	      	callback(response)
	      })
	      .catch(function (error) {
	    	console.log(error);
	  	  })
	  	  .then(function () {
	  	   });
	}
	delete(id,callback){
		axios.delete(`${HOST}/api/clients/${id}`)
	      .then((response) => {
	      	callback(response)
	      })
	      .catch(function (error) {
	    	console.log(error);
	  	  })
	  	  .then(function () {
	  	   });
	}

}

export { 
	ApiClient,
	ApiChannel,
	MessageContent,
	ApiMessage
}