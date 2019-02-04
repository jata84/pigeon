<template>

      <v-container >
          <v-card>
            <v-card-title primary-title>
              <div>
                <h3 class="headline mb-0">Channel: {{ channel.name }}</h3>
              </div>
            </v-card-title>
            <v-card-text>
              <v-textarea
                      v-model="message"
                      name="input-7-4"
                      label="Message"
                      value=""
              ></v-textarea>
            </v-card-text>

            <v-card-actions>
              <v-btn flat color="primary" v-on:click="sendMessage">Enviar</v-btn>
            </v-card-actions>
          </v-card>

      </v-container>



</template>
<script>

import { RTApiClient } from '@/lib/websocket_api.js'
import { MessageContent } from '@/lib/api.js'
import { ApiMessage } from '@/lib/api.js'


  export default {
    data () {
      return {
        message:"",
      	test:"aaa",
        channel:{},
        channel_clients:[]
      }
    },
    mounted(){
      let that = this

      let name = this.$route.params.name

      this.RTApiClient.retrieve(name,function(response){
        that.channel = data
      })

      this.RTApiClient.getClients(name,function(response){
        that.channel_clients = response.data
      })
      
    },
     methods:{

      sendMessage: function(event){
        this.RTApiClient.message(this.message,[this.channel.name],(data) =>{
          console.log(data)
        })
        this.message=""

      }
    }
  }
</script>