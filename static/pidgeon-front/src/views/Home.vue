<template>
 <v-container fluid fill-height>
        <v-layout
          justify-center
          align-center
        >
          <v-flex text-xs-center>
			  <div class="about">
			    <h1>This is an about page</h1>
			  </div>
          </v-flex>
        </v-layout>
      </v-container>

</template>


<script>
import Home from '@/components/home.vue'

export default {
  name: 'home',
  components: {
    Home
  }
}
</script>
