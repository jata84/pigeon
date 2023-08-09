import Vue from 'vue'
import Vuex from 'vuex'

Vue.use(Vuex)


const connections = {
  namespaced: true,
  state: {
    connections: [],
  },
  mutations: { 
    addConnection (state,connection){
      state.connections.push(connection)
      },
    updateConnection (state,id,connection){
              let connectionElement = state.connections.findIndex(function (element) {
                return element.idClient === id
              })
              if (connectionElement>-1){
                state.connections[connectionElement].messageList = connection.messageList
              }
              
      },
    deleteConnection (state,id,connection){
              let connectionElement = state.connections.findIndex(function (element) {
                return element.idClient === id
              })
              if (connectionElement>-1){
                state.connections.splice(connectionElement,1)
              }
              
      }
  },
  actions: {

  },
  getters: {
      getConnection: (state,id) => {
            let connectionElement = state.connections.findIndex(function (element) {
                return element.idClient === id
            })
        },

      getConnections: state => {
          return state.connections
      }
  }
}

const clients = {
  state: {
    clients: {},
  },
  mutations: {
    addClient (state,client){
      state.clients[client.uuid] = client
      },

    addClients (state,clients){
          clients.map((c)=>{
              state.clients[c.uuid]=c
          })
          //state.clients[client.uuid] = state.clients.concat(clients)
      },
    updateClient (state,client){
              /*let connectionElement = state.clients.findIndex(function (element) {
                return element.idClient === id
              })
              if (connectionElement>-1){
                state.clients[connectionElement].messageList = client.messageList
              }*/
              state.clients[client.uuid]=client

              
      },
    deleteClient (state,id){
        /*
              let connectionElement = state.clients.findIndex((element) =>{
                return element.uuid === id
              })
        console.log(connectionElement)
              if (connectionElement>-1){
                state.clients.splice(connectionElement,1)
              }
         */
        delete state.clients[id];

      }

  },
  actions: {

  },
  getters: {
      getClients: state => {
          let arr = [];
          for (let key in state.clients) {
              if (state.clients.hasOwnProperty(key)) {
                  arr.push( [ key, state.clients[key] ] );
              }
          }

          return arr
      }
  }
}



const channels = {
  state: {
    channels: [],
  },
  mutations: {
    addChannels (state,channel){
        state.channels.push(channel)
      },
    updateChannels (state,name,channel){
              let connectionElement = state.channels.findIndex(function (element) {
                return element.name === id
              })
              if (connectionElement>-1){
              
              }
              
      },
    deleteChannels (state,name,channel){
              let connectionElement = state.channels.findIndex(function (element) {
                return element.name === name
              })
              if (connectionElement>-1){
                state.channels.splice(connectionElement,1)
              }
              
      }

  },
  actions: {

  },
  getters: {

  }
}



export default new Vuex.Store({
  modules:{
    connections,
    clients,
    channels
  }
})

