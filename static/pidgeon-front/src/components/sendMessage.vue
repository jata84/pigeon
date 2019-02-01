<template>
  <v-container >
          <v-card>
            <v-card-text>
            	  <v-combobox
			          v-model="select"
			          :items="channel_list"
			          item-text="name"
			          chips
			          multiple
			          label="Select Channels"
                      v-on:select="refreshNamespace"
			        ></v-combobox>
              <v-textarea
                      v-model="message"
                      name="input-7-4"
                      label="Message"
                      value=""
              ></v-textarea>
            </v-card-text>

            <v-card-actions>
              <v-btn flat color="primary" v-on:click="sendMessage">Send</v-btn>
            </v-card-actions>
          </v-card>
  </v-container>
</template>

<script>

import { ApiChannel } from '@/lib/api.js'
import { MessageContent } from '@/lib/api.js'


  export default {
    data () {
      return {
		select: '',
		channel_list: [
		        
		      ],
        message:"",
      	test:"aaa",
        channel:{},
        channel_clients:[]
      }
    },
    mounted(){
      let that = this
      let name = this.$route.params.name
      this.RTApiClient.list_namespace(function(data){
        that.channel_list = data
      })
      
    },
     methods:{
      refreshNamespace:function(event){
          let that = this
          this.RTApiClient.list_namespace(function(data){
              that.channel_list = data
          })
      },

      sendMessage: function(event){
        this.RTApiClient.message(this.message,this.select.map(function(d){return d.name}),(response) =>{
            console.log(response.data)
        })
        this.message=""

      }
    }
  }
</script>

<style lang="css" scoped>
</style>