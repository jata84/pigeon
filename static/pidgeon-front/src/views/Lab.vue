<template>
	<v-container>
		<v-layout row wrap>    
      		<v-flex xs12 sm6>
					<SendMessage></SendMessage>

      			
      		</v-flex>
			<v-flex xs12 sm6 >
					<CreateConnection></CreateConnection> 
			      	
      		</v-flex>
      	</v-layout>

		
	</v-container>
	
</template>
<script>

import SendMessage from '@/components/sendMessage.vue'
import CreateConnection from '@/components/createConnections.vue'

export default {

  name: 'labview',

  data () {
    return {

    }
  },    
  computed: { 
      connections: function(){
        return this.$store.state.connections.getters.getConnections 
      }
  },
  components:{
      SendMessage,
      CreateConnection
    }

}
</script>

<style lang="css" scoped>
</style>

