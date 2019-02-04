import Vue from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
// index.js or main.js
import 'vuetify/dist/vuetify.min.css' // Ensure you are using css-loader
import vuetify from 'vuetify'
import {RTAuthentication,RTApiClient,BASIC} from '@/lib/websocket_api.js'

let authentication = new RTAuthentication(BASIC,"test")
console.log(authentication.getToken())
Vue.prototype.RTApiClient = new RTApiClient("localhost","8084",authentication.getToken());

Vue.use(vuetify)
Vue.config.productionTip = false



new Vue({
  router,
  store,
  render: h => h(App)
}).$mount('#app')
