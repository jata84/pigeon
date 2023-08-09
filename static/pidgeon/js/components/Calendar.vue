<template>
  <v-container
      fluid
      grid-list-lg
      class="mt-0"
    >
    <v-layout row wrap>
      <v-flex xs12 md6 v-for="item in list" :key="item.id">
        <v-card tile>
          <v-card-media
            class="white--text"
            :src="item.imageURL" 
            height="200px">
            <v-container fluid class="mb-0 pl-2 pb-2">
              <v-layout>
                <v-flex xs12>
                  <v-icon x-large dark>event</v-icon>
                  <span class="subheading">{{ item.date }} {{ item.hour }}</span>
                </v-flex>
              </v-layout>
            </v-container>
          </v-card-media>
          <v-card-title primary-title>
            <div>
              <h3 class="headline">{{ item.title }}</h3>
              <div class="subheading pl-0">
                <v-icon>location_on</v-icon>&nbsp;{{ item.location }}
              </div>
              <v-flex xs12 class="grey--text text--darken-1 pl-0">
                <v-icon class="grey--text text--lighten-1">supervisor_account</v-icon>
                <span class="body-1">&nbsp;{{ item.joined }} People Joined</span>
              </v-flex>
            </div>
          </v-card-title>
          <v-card-actions>
            <v-btn flat>Share</v-btn>
            <v-btn flat
              :to="
                  {
                    name:'CalendarDetail', 
                    params:
                    {
                      id:item.id, 
                      item:item
                    }
                  }"
              >
              View Detail
            </v-btn>
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
      data:() => ({
        list: [
          {
            id: 1,
            title: "Ceremony at Monkey Forest",
            location: "Alas Kedaton, Tabanan",
            date: "18 Mar 2018",
            hour: "14:00 PM",
            joined: 32,
            imageURL: "https://c2.staticflickr.com/8/7390/11206634984_44f61f926d.jpg",
            desc: "The temple ceremony in Alas Kedaton Temple is carried out every 210 days a year. It is on <i>Anggarakasih Medangsia</i> (Balinese Hindu calendar) or on every Tuesday where on that time the society do the worship or pray to request the safety and prosperity. The unique in this ceremony is do not use the fire and do not hence Penjor and also finished before the sunset or before the night is come."
          },
          {
            id: 2,
            title: "Pandanus War",
            location: "Tenganan Pegeringsingan, Karangasem",
            date: "20 Jul 2018",
            hour: "10:00 AM",
            joined: 40,
            imageURL: "https://c2.staticflickr.com/6/5331/17634071142_152c806768.jpg",
            desc: "Pandanus War or localy known as <i>Mekare-kare</i> is an annual theatrical fight between the young men of the village, utilising prickly pandanus leaf whips! Each dual is staged to the intense martial sounds of <i>gamelan selonding</i> music, and lasts only a few seconds, accompanied by much merriment and laughter."
          },
          {
            id: 3,
            title: "Bali Mass Cremation",
            location: "Tampaksiring, Gianyar",
            date: "6 Aug 2018",
            hour: "09:00 AM",
            joined: 10,
            imageURL: "https://c1.staticflickr.com/9/8713/17470094589_0f1fd7c67f.jpg",
            desc: "Cremation is guiding the soul of the deceased family member from hell and underworld to divine realms, later to be reborn in the next generation of the same family. Long preparations proceed cremation ceremonies.<br>Each family uses a certain type of sarcophagus, sculptures formed as mythical animals, according to their cast. A black or white bull is used for a Brahman cast, together with a tower-shaped sarcophage called bade, with roofs like a pagode."
          },
          {
            id: 4,
            title: "Makepung Lampit",
            location: "Yeh Peh, Jembrana",
            date: "18 Nov 2018",
            hour: "07:00 AM",
            joined: 24,
            imageURL: "https://c1.staticflickr.com/9/8117/8678393979_1997fecfe8.jpg",
            desc: "Makepung Lampit is a kind of buffalo race without using the wheel as often held now in Jembrana district of Bali, but using the thing called <i>Lampit</i>. Lampit is a tool to flatten rice fields before planting. <br>Makepung Lampit have started to be implemented since the 1940 and stopped around 1990 after the available machines such as tractors to plow rice fields."
          }
        ]
      })/*,
      mounted(){
        console.log(this.$route);
      }*/
    });
  });
</script>
<style>
  .joined {
    /*color: red;*/
  }
</style>