define(["VueRouter", 
        "vue!components/Home",
        "vue!components/Calendar",
        "vue!components/CalendarDetail",
        "vue!components/Profile",
        "vue!components/Notification"
    ],
    function(VueRouter, 
        Home,
        Calendar,
        CalendarDetail,
        Profile,
        Notification
    )
{
  Vue.use(VueRouter);
  
  return new VueRouter({
    routes: [
      {
          path: "/calendar", 
          name: "Calendar", 
          component:Calendar
      },
      {
          path: "/calendar/:id", 
          name: "CalendarDetail", 
          component:CalendarDetail,
          props: true
      },
      {
          path: "/profile", 
          name: "Profile",
          component:Profile
      },
      {
          path: "/notification", 
          name: "Notification",
          component:Notification
      },
      {
          path: "/", 
          name: "Home", 
          component:Home
      }
    ],
    //mode: 'history'
  });
});