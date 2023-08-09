<template>
  <v-container
      fluid
      grid-list-lg
      class="mt-0"
    >
    <v-layout row wrap>
      <v-flex xs12>
        <v-card tile>
          <v-carousel hide-delimiters class="mycarousel">
              <v-carousel-item dark v-for="item in list" :key="item.id" :src="item.imageURL">
                <div class="white--text text-xs-center">
                  <h3 class="headline">{{ item.title }}</h3>
                  <div class="subheading">
                    <v-icon dark>location_on</v-icon>&nbsp;{{ item.location }}
                  </div>
                  <div class=""title>
                    {{ item.date }}  {{ item.hour }}
                  </div>
                </div>
              </v-carousel-item>
          </v-carousel>
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
      data:() => ({
        list: [
          {
            id: 1,
            title: "Ceremony at Monkey Forest",
            location: "Alas Kedaton, Tabanan",
            date: "18 Mar 2018",
            hour: "14:00 PM",
            joined: 32,
            imageURL: "https://c2.staticflickr.com/8/7390/11206634984_44f61f926d.jpg"
          },
          {
            id: 2,
            title: "Pandanus War",
            location: "Tenganan Pegeringsingan, Karangasem",
            date: "20 Jul 2018",
            hour: "10:00 AM",
            joined: 40,
            imageURL: "https://c2.staticflickr.com/6/5331/17634071142_152c806768.jpg",
          },
          {
            id: 3,
            title: "Bali Mass Cremation",
            location: "Tampaksiring, Gianyar",
            date: "6 Aug 2018",
            hour: "09:00 AM",
            joined: 10,
            imageURL: "https://c1.staticflickr.com/9/8713/17470094589_0f1fd7c67f.jpg"
          },
          {
            id: 4,
            title: "Makepung Lampit",
            location: "Yeh Peh, Jembrana",
            date: "18 Nov 2018",
            hour: "07:00 AM",
            joined: 24,
            imageURL: "https://c1.staticflickr.com/9/8117/8678393979_1997fecfe8.jpg"
          }
        ]
      })
    });
  });
</script>