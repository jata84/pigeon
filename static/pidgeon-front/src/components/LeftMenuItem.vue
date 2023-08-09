<template>
	  <v-list-tile @click="">
          <v-list-tile-action>
            <v-icon>{{ iconName }}</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-list-tile-title><slot></slot></v-list-tile-title>
          </v-list-tile-content>
        </v-list-tile>

</template>

<script>
export default {
	props: ['iconName'],
}
</script>
