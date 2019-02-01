<template>
  <v-container
      fluid
      grid-list-lg
      class="mt-0"
    >
    <v-layout row wrap>
      <v-flex xs12>
        <v-card tile>
          <!--v-card-media :src="item.imageURL" contain>
          </v-card-media-->
          <v-carousel hide-delimiters class="mycarousel">
              <v-carousel-item :src="item.imageURL"></v-carousel-item>
          </v-carousel>
          <v-card-title primary-title>
            <div>
              <h3 class="headline">{{ item.title }}</h3>
              <div class="subheading pl-0">
                <v-icon>location_on</v-icon>&nbsp;{{ item.location }}
              </div>
              <v-flex xs12 class="pl-0">
                  <v-icon>event</v-icon>
                  <span class="body-1">{{ item.date }}  {{ item.hour }}</span>
              </v-flex>
              <v-flex xs12 class="grey--text text--darken-1 pl-0">
                <v-icon class="grey--text text--lighten-1">supervisor_account</v-icon>
                <span class="body-1">&nbsp;{{ item.joined }} People Joined</span>
              </v-flex>
              <div v-if="item.desc" v-html="item.desc"></div>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-btn flat>Share</v-btn>
            <v-btn flat to="/calendar">Close</v-btn>
          </v-card-actions>
        </v-card>
      </v-flex>
    </v-layout>
  </v-container>
</template>
<script>
 define([], function()
 {
    return Vue.extend({
        template,
        props: ["id", "item"]
        /*data:() => (){
            
        })*/
    });
 });
</script>
<style>
@media only screen and (max-width: 320px) {  
  .mycarousel {
    height:100px;
  }
}
@media only screen and (max-width: 600px) {  
  .mycarousel {
    height:250px;
  }
}
@media only screen and (max-width: 961px) {  
  .mycarousel {
    height:350px;
  }
}
</style>