package main

import (
	"net/http"

	"github.com/labstack/echo"
)

// https://66da3a82c5f1.ngrok.io/SiprtcApplications/MainRestaurantMenu

// var input = []byte(`{"text":"i want to book a table for 4","intent":{"name":"booking_with_count","confidence":0.9603162407875061},"entities":[{"start":27,"end":28,"text":"4","value":4,"confidence":1,"entity":"number"}],"intent_ranking":[{"name":"booking_with_count","confidence":0.9603162407875061},{"name":"booking","confidence":0.0258566252887249},{"name":"booking_with_count_time_day_hours_minute","confidence":0.0034164865501224995},{"name":"booking_with_count_time","confidence":0.0032817136961966753},{"name":"cancel_booking","confidence":0.00111346784979105},{"name":"goodbye","confidence":0.0010961686493828893},{"name":"order_pizza","confidence":0.0008867786964401603},{"name":"bot_challenge","confidence":0.0008118133991956711},{"name":"complain","confidence":0.000798406545072794},{"name":"booking_time_day","confidence":0.0006561993504874408}]}`)

// var base64Input = `eyJ0ZXh0IjoiaSB3YW50IHRvIGJvb2sgYSB0YWJsZSBmb3IgNCIsImludGVudCI6eyJuYW1lIjoiYm9va2luZ193aXRoX2NvdW50IiwiY29uZmlkZW5jZSI6MC45NjAzMTYyNDA3ODc1MDYxfSwiZW50aXRpZXMiOlt7InN0YXJ0IjoyNywiZW5kIjoyOCwidGV4dCI6IjQiLCJ2YWx1ZSI6NCwiY29uZmlkZW5jZSI6MSwiZW50aXR5IjoibnVtYmVyIn1dLCJpbnRlbnRfcmFua2luZyI6W3sibmFtZSI6ImJvb2tpbmdfd2l0aF9jb3VudCIsImNvbmZpZGVuY2UiOjAuOTYwMzE2MjQwNzg3NTA2MX0seyJuYW1lIjoiYm9va2luZyIsImNvbmZpZGVuY2UiOjAuMDI1ODU2NjI1Mjg4NzI0OX0seyJuYW1lIjoiYm9va2luZ193aXRoX2NvdW50X3RpbWVfZGF5X2hvdXJzX21pbnV0ZSIsImNvbmZpZGVuY2UiOjAuMDAzNDE2NDg2NTUwMTIyNDk5NX0seyJuYW1lIjoiYm9va2luZ193aXRoX2NvdW50X3RpbWUiLCJjb25maWRlbmNlIjowLjAwMzI4MTcxMzY5NjE5NjY3NTN9LHsibmFtZSI6ImNhbmNlbF9ib29raW5nIiwiY29uZmlkZW5jZSI6MC4wMDExMTM0Njc4NDk3OTEwNX0seyJuYW1lIjoiZ29vZGJ5ZSIsImNvbmZpZGVuY2UiOjAuMDAxMDk2MTY4NjQ5MzgyODg5M30seyJuYW1lIjoib3JkZXJfcGl6emEiLCJjb25maWRlbmNlIjowLjAwMDg4Njc3ODY5NjQ0MDE2MDN9LHsibmFtZSI6ImJvdF9jaGFsbGVuZ2UiLCJjb25maWRlbmNlIjowLjAwMDgxMTgxMzM5OTE5NTY3MTF9LHsibmFtZSI6ImNvbXBsYWluIiwiY29uZmlkZW5jZSI6MC4wMDA3OTg0MDY1NDUwNzI3OTR9LHsibmFtZSI6ImJvb2tpbmdfdGltZV9kYXkiLCJjb25maWRlbmNlIjowLjAwMDY1NjE5OTM1MDQ4NzQ0MDh9XX0=`

var numberRestMap *PhonenumberMap

func main() {

	e := echo.New()

	numberRestMap = new(PhonenumberMap)

	e.GET("/v1/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy!!!")
	})

	// Voxvelly basic DTMF and Call fowarding demo
	e.POST("/SiprtcApplications/VoxvellyDemo", VoxvellyDemo)
	e.POST("/SiprtcApplications/VoxvellyDemoDtmfReceived", VoxvellyDemoDtmfReceived)

	// Directcall application demo
	e.POST("/SiprtcApplications/DirectCall", DirectCall)

	// Restaurent demo DTMF bot
	e.GET("/SiprtcApplications/MainRestaurantMenu", MainRestaurantMenu)
	e.POST("/SiprtcApplications/DtmfReceived", RestaurantDtmfReceived)

	// Restaurent demo voice bot
	e.POST("/SiprtcApplications/Voicebot", VoicebotWelcomeMessage)
	e.POST("/SiprtcApplications/UserIntent", VoicebotUserIntent)

	// ProcessApplicationWelcomeXML
	e.GET("/SiprtcApplications/ProcessApplicationWelcomeXML", ProcessApplicationWelcomeXML)

	e.Logger.Fatal(e.Start(":8080"))
}
