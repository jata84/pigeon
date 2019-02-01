
<template>
  <v-app id="inspire">
    <v-navigation-drawer
      v-model="drawer"
      fixed
      app
    >
      <v-list dense>
        <LeftMenuItem icon-name="home"><router-link to="/home">Home</router-link></LeftMenuItem>
        <LeftMenuItem icon-name="contact_mail"><router-link to="/about">About</router-link></LeftMenuItem>
        <LeftMenuItem icon-name="forum"><router-link to="/namespace">Channels</router-link></LeftMenuItem>
        <LeftMenuItem icon-name="compare_arrows"><router-link to="/connections">Connections</router-link></LeftMenuItem>
        <LeftMenuItem icon-name="extension"><router-link to="/playground">PlayGround</router-link></LeftMenuItem>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar color="indigo" dark fixed app>
      <v-toolbar-side-icon @click.stop="drawer = !drawer"></v-toolbar-side-icon>
      <v-toolbar-title>Application</v-toolbar-title>
    </v-toolbar>
    <v-content>
    <router-view></router-view>
     
    </v-content>
    <v-footer color="indigo" app>
      <span class="white--text">&copy; 2017</span>
    </v-footer>
  </v-app>
</template>

<script>

import LeftMenuItem from '@/components/LeftMenuItem.vue'

  export default {
    created: function () {
    },
    data: () => ({
      connections: null,
      drawer: null
    }),
    props: {
      source: String
    },
    components:{
      LeftMenuItem
    }
  }
</script>


<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
