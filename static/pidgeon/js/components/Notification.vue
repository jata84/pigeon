<template>
  <v-container class="ma-0 pa-10">
    <v-alert color="warning" icon="info" value="true">
      Today is the last opportunity to join Pandanus War
    </v-alert>
    <v-alert color="error" icon="warning" value="true">
      Sorry you've missed Barong Brutuk Sacred Dance yesterday
    </v-alert>
    <v-alert color="success" icon="check_circle" value="true">
      You joined Ceremony at Monkey Forest yesterday.
    </v-alert>
    <v-alert color="info" icon="info" value="true">
      Congrats...!! You registered Photo Hunt Calendar a year ago
    </v-alert>
  </v-container>
</template>
<script>
    define([], function(){
        return Vue.extend({
            template
        });
    });
</script>