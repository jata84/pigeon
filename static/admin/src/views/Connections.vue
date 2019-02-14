<template>
    <v-container >
        <v-data-table
                :headers="headers"
                :items="client_list"
                class="elevation-1"
                hide-actions
        >
            <template slot="items" slot-scope="props">
                <td>{{ props.item.uuid }}</td>
                <td class="text-xs-right">{{ props.item.datetime }}</td>
                <td class="text-xs-right">
                    <v-btn color="blue darken-1" flat @click="close(props.item)">Close</v-btn>
                </td>
            </template>
        </v-data-table>
    </v-container>
</template>
<script>
import {RTApiClient} from '@/lib/websocket_api.js'
import { ApiClient } from '@/lib/api.js'
import { mapGetters } from 'vuex';
  export default {
    data () {
      return{
          headers: [
              {
                  text: 'uuid',
                  value: 'uuid'
              },
              {
                  text: 'date',
                  value: 'datetime'
              },
              {
                  text: 'action',
                  value: 'action'
              },
          ],

      }
    },

    mounted(){
      let that = this
      this.RTApiClient.list_connections((data)=>{
          this.$store.commit('addClients',data)
      })
    },
    computed: {

        /*
        ...mapGetters({
            client_list: 'getClients',
        }),
        */
    },
    methods:{
        close: function(item){
            this.RTApiClient.close(item.uuid,(data)=>{
                this.$store.commit('deleteClient', item.uuid)

            })
        },
    },
  }
</script>