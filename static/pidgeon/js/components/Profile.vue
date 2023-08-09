<template>
    <v-container class="mt-0">
    {{ title }}
    </v-container>
</template>
<script>
  define([], function()
  {
    return Vue.extend({
      template,
      data:() => ({
        title: "Profile Info"
      })
    });
  });
</script>