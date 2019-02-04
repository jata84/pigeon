(function()
{
  require.config(
  {
    baseUrl: "js/",
    paths:{
      vue: "https://cdn.rawgit.com/edgardleal/require-vuejs/aeaff6db/dist/require-vuejs.min",
      Vuetify: "https://unpkg.com/vuetify/dist/vuetify",
      VueRouter: "https://unpkg.com/vue-router/dist/vue-router"//"https://cdnjs.cloudflare.com/ajax/libs/vue-router/2.1.1/vue-router.min"
    }
  });
  
  require(["Vuetify", "router/router", "vue!App"], 
          function(Vuetify, router, App)
  {
    //Vue.use(Vuetify);
    
    Vue.config.productionTip = false;
    
    new Vue({
      el: '#app',
      router: router,
      render: h => h(App)
    })
    
  });
  
})();

