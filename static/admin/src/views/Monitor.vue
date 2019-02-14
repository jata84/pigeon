<template>
    <v-container >
        <span>Multiline message is:</span>
        <div>
            <p> {{ messages }}</p>
        </div>

    </v-container>


</template>

<script>
    export default {
        name: "Monitor",
        data () {
            return {
                messages:"test"
            }
        },
        methods:{
            close: function(item){
                this.RTApiClient.close(item.uuid,(data)=>{
                    this.$store.commit('deleteClient', item.uuid)

                })
            },
        }
    }
</script>

<style scoped>

</style>