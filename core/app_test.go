package core

/*
func TestBasicApp(t *testing.T) {
	LoadConfig()
	app := NewApp()
	go app.Init()
	time.Sleep(4000 * time.Millisecond)
	client := client.NewWebsocketClient("localhost", "4444")
	client.Init()
	client.SendRegistration()

	resp, _ := grequests.Post("http://localhost:8084/api/message/",
		&grequests.RequestOptions{JSON: []byte(
			`{
    			"timestamp":"2018-10-17T19:51:54.126486939+02:00",
    			"type":"MESSAGE",
    			"comunication":{
      			"type":"MESSAGE",
      			"namespaces":[{"name":"test"}]
			},
    		"data":{"data":"ttttt"}
 			}`)})

	fmt.Printf("Post: %v", resp)

	client.SendSubscription("test")


	time.Sleep(3000 * time.Millisecond)

}
*/
