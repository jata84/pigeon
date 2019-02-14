<template>
	<v-container >
		<v-card>
			<v-btn color="primary" v-on:click="create_connection">Create Connection</v-btn>
			<v-card v-for="item in connections"></v-card>
			<Connection v-for="item in connections" :connection= "item"> </Connection>
		</v-card>
	</v-container>
</template>

<script>
	import {RTClient} from '@/lib/websocket_api.js'
	import {ApiChannel} from '@/lib/api.js'
	import Connection from '@/components/connection.vue'
	import Store from '@/store.js'
	export default {
		name:'createConnections',
		data () {
			return{

			}
		},
		mounted(){
		},
		computed: {
			connections: function(){
				return this.$store.getters['connections/getConnections']
			}
		},
		methods:{
			create_connection: function(event){
				var client = new RTClient("localhost"
						,"4444"
						,(data)=>{
							client.messageList.push(data)
							Store.commit('connections/updateConnection',client)
						}
						,()=>{
							Store.commit('connections/deleteConnection',client.idClient,client)
						}
						,()=>{

							Store.commit('connections/addConnection',client)
						},()=>{
							Store.commit('connections/updateConnection',client)
						}
				)

			}
		},
		components:{
			Connection
		}


	}
</script>

<style lang="css" scoped>
</style>
