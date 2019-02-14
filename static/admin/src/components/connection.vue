<template>
    <v-container  grid-list-md>
        <v-card>
            <v-layout row wrap>
                <v-flex xs12>
                      <v-card-text class="px-12"> <b>Connection Id:</b> {{ connection.idClient }}</v-card-text>
                </v-flex>
                <v-flex xs12>
                    <div class="d-inline-flex" v-for="(item, index) in connection.namespaces_subscribed">
                        <v-chip close @input="unsubscribe(item)">{{ item }}</v-chip>
                    </div>
                </v-flex>
                <v-flex xs12>
                      <v-card-text v-if="connection.messageList.length > 0" class="px-5"><b>Last Message received:</b> {{ connection.messageList[connection.messageList.length-1].content }}</v-card-text>
                </v-flex>
            </v-layout>
            <v-card-actions>
              <v-btn v-on:click="close">Close</v-btn>
                <v-dialog v-model="dialog" persistent max-width="600px">
                    <v-btn slot="activator" color="primary" dark>Subscription</v-btn>
                    <v-card>
                        <v-card-title>
                            <span class="headline">Connection Details</span>
                        </v-card-title>
                        <v-card-text>
                            <v-container grid-list-md>
                                <v-layout wrap>

                                    <v-flex xs12>
                                        <b>id:{{connection.idClient}}</b>
                                        <v-text-field v-model="namespace" label="Namespace" required></v-text-field>
                                        <v-btn v-on:click="subscribe(namespace)" color="primary" dark>Subscribe</v-btn>
                                    </v-flex>


                                </v-layout>
                            </v-container>
                        </v-card-text>
                        <v-card-actions>
                            <v-spacer></v-spacer>
                            <v-btn color="blue darken-1" flat @click="dialog = false">Close</v-btn>
                        </v-card-actions>
                    </v-card>
                </v-dialog>
            </v-card-actions>
        </v-card>
    </v-container>
</template>
<script>
  export default {
    data () {
      return{
          dialog: false,
          namespace:"",
      }
    },
    props: {
      connection:{},
    },
    computed: {
    },

    methods:{

          unsubscribe:function(element){
              this.RTApiClient.unsubscription([element],this.connection.idClient,()=>{
                  /*let index = this.connection.namespaces_subscribed.indexOf(element)
                  if (index > -1){
                      this.connection.namespaces_subscribed.splice(index,1)
                  }*/
              })
          },
          close: function(){
              this.connection.close()
          },
          detail: function(){

          },
          subscribe:function(namespace){
              this.RTApiClient.subscription([namespace],this.connection.idClient)
              namespace=""
          }

      }

  }
</script>