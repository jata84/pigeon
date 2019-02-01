import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import Channels from './views/Channels.vue'
import ChannelDetail from './views/ChannelDetail.vue'
import About from './views/About.vue'
import Connections from './views/Connections.vue'
import Lab from './views/Lab.vue'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/namespace',
      name: 'namespace',
      component: Channels
    },
    {
      path: '/namespace/:name',
      name: 'namespace/detail',
      component: ChannelDetail,
      props: true
    },
    {
      path: '/about',
      name: 'about',
      component: About
    },
    {
      path: '/connections',
      name: 'connections',
      component: Connections
    },
    {
      path: '/playground',
      name: 'playground',
      component: Lab
    }
  ]
})
