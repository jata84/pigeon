<template>
    <v-container  grid-list-md>
        <v-card>
            <v-layout row wrap>
                <v-flex xs6>
                      <v-card-text class="px-5"> <b>Connection Id:</b> {{ client.uuid }}</v-card-text>
                </v-flex>
            </v-layout>
            <v-card-actions>
              <v-btn v-on:click="close">Close</v-btn>
            </v-card-actions>
        </v-card>
    </v-container>
</template>
<script>
  export default {
    data () {
      return{
        
      }
    },
    props: {
      client:{}
    },
    mounted(){      
    },
    computed: { 
    },
    methods:{


    }
  }
</script>