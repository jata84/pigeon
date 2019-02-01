<template>
  <v-container>

  <v-data-table
    :headers="headers"
    :items="namespaces"
    hide-actions
    class="elevation-1"
  >
    <template slot="items" slot-scope="props">
      <td>{{ props.item.name }}</td>
      <td class="text-xs-left">{{ props.item.type }}</td>
      <td class="text-xs-left">
        <router-link :to="{ name: 'namespace/detail', params: { name: props.item.name }}"> Detalle</router-link>
      </td>
    </template>
  </v-data-table>


  </v-container>

</template>

<script>

import {ApiChannel} from '@/lib/api.js'
import {RTApiClient} from '@/lib/websocket_api.js'
  export default {
    data () {
      return {
        headers: [
          {
            text: 'Namespace name',
            align: 'left',
            sortable: false,
            value: 'name'
          },
          { text: 'Type', value: 'type' },
          { text: 'Name', value: 'name' },


        ],
        namespaces: [
        
        ]
      }
    },
    created: function () {

      
    },
    mounted(){
      let that = this
      this.RTApiClient.list_namespace(function(data){
        that.namespaces = data
      })
    }
  }
</script>