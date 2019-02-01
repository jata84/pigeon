<template>
  <v-app>
    <v-navigation-drawer 
      absolute
      class="hidden-sm-and-up"
      v-model="sideNav">
      <v-toolbar flat>
        <v-list>
          <v-list-tile>
            <v-list-tile-title class="title">
              {{ appName }}
            </v-list-tile-title>
          </v-list-tile>
        </v-list>
      </v-toolbar>
      <v-list>
        <v-list-tile 
          v-for="item in menuItems" 
          :key="item.title"
          :to="item.link"
          >
          <v-list-tile-action>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-tile-action>
          <v-badge color="red" v-if="item.badge">
            <span slot="badge">{{ item.badge }}</span>
            <span>{{ item.title }}</span>
          </v-badge>
          <v-list-tile-content v-else>{{ item.title }}</v-list-tile-content>
        </v-list-tile>
      </v-list>
    </v-navigation-drawer>
    <v-toolbar>
      <v-toolbar-side-icon 
        @click.stop="sideNav = !sideNav"
        class="hidden-sm-and-up"></v-toolbar-side-icon>
      <v-toolbar-title>
        <router-link 
          to="/"
          tag="span"
          style="cursor: pointer"
          >
          {{ appName }}
        </router-link>
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-toolbar-items class="hidden-xs-only">
        <v-btn flat 
          v-for="item in menuItems" 
          :key="item.title"
          :to="item.link"
          >
          <v-icon left>{{ item.icon }}</v-icon>
          <v-badge color="red" v-if="item.badge">
            <span slot="badge">{{ item.badge }}</span>
            <span>{{ item.title }}</span>
          </v-badge>
          <span v-else>{{ item.title }}</span>
        </v-btn>
      </v-toolbar-items>  
    </v-toolbar>
    <main>
      <router-view></router-view>
    </main>  
  </v-app>
</template>
<script>
  define([], function()
  {
    return Vue.extend(
    {
      template: template,
      data: function() {
        return {
            sideNav: false,
            appName: "Photo Hunt Calendar",
            menuItems: [
                {icon: 'date_range', title: 'Event Calendar', link: '/calendar'},
                //{icon: 'person', title: 'Profile', link: '/profile'},
                {icon: 'notifications', title: 'Notifications', badge:4, link: '/notification'}
                
            ]
        }
      }
    });
  });
</script>