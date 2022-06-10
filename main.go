package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

// https://66da3a82c5f1.ngrok.io/SiprtcApplications/MainRestaurantMenu

// var input = []byte(`{"text":"i want to book a table for 4","intent":{"name":"booking_with_count","confidence":0.9603162407875061},"entities":[{"start":27,"end":28,"text":"4","value":4,"confidence":1,"entity":"number"}],"intent_ranking":[{"name":"booking_with_count","confidence":0.9603162407875061},{"name":"booking","confidence":0.0258566252887249},{"name":"booking_with_count_time_day_hours_minute","confidence":0.0034164865501224995},{"name":"booking_with_count_time","confidence":0.0032817136961966753},{"name":"cancel_booking","confidence":0.00111346784979105},{"name":"goodbye","confidence":0.0010961686493828893},{"name":"order_pizza","confidence":0.0008867786964401603},{"name":"bot_challenge","confidence":0.0008118133991956711},{"name":"complain","confidence":0.000798406545072794},{"name":"booking_time_day","confidence":0.0006561993504874408}]}`)

var base64Input = `eyJ0ZXh0IjoiaSB3YW50IHRvIGJvb2sgYSB0YWJsZSBmb3IgNCIsImludGVudCI6eyJuYW1lIjoiYm9va2luZ193aXRoX2NvdW50IiwiY29uZmlkZW5jZSI6MC45NjAzMTYyNDA3ODc1MDYxfSwiZW50aXRpZXMiOlt7InN0YXJ0IjoyNywiZW5kIjoyOCwidGV4dCI6IjQiLCJ2YWx1ZSI6NCwiY29uZmlkZW5jZSI6MSwiZW50aXR5IjoibnVtYmVyIn1dLCJpbnRlbnRfcmFua2luZyI6W3sibmFtZSI6ImJvb2tpbmdfd2l0aF9jb3VudCIsImNvbmZpZGVuY2UiOjAuOTYwMzE2MjQwNzg3NTA2MX0seyJuYW1lIjoiYm9va2luZyIsImNvbmZpZGVuY2UiOjAuMDI1ODU2NjI1Mjg4NzI0OX0seyJuYW1lIjoiYm9va2luZ193aXRoX2NvdW50X3RpbWVfZGF5X2hvdXJzX21pbnV0ZSIsImNvbmZpZGVuY2UiOjAuMDAzNDE2NDg2NTUwMTIyNDk5NX0seyJuYW1lIjoiYm9va2luZ193aXRoX2NvdW50X3RpbWUiLCJjb25maWRlbmNlIjowLjAwMzI4MTcxMzY5NjE5NjY3NTN9LHsibmFtZSI6ImNhbmNlbF9ib29raW5nIiwiY29uZmlkZW5jZSI6MC4wMDExMTM0Njc4NDk3OTEwNX0seyJuYW1lIjoiZ29vZGJ5ZSIsImNvbmZpZGVuY2UiOjAuMDAxMDk2MTY4NjQ5MzgyODg5M30seyJuYW1lIjoib3JkZXJfcGl6emEiLCJjb25maWRlbmNlIjowLjAwMDg4Njc3ODY5NjQ0MDE2MDN9LHsibmFtZSI6ImJvdF9jaGFsbGVuZ2UiLCJjb25maWRlbmNlIjowLjAwMDgxMTgxMzM5OTE5NTY3MTF9LHsibmFtZSI6ImNvbXBsYWluIiwiY29uZmlkZW5jZSI6MC4wMDA3OTg0MDY1NDUwNzI3OTR9LHsibmFtZSI6ImJvb2tpbmdfdGltZV9kYXkiLCJjb25maWRlbmNlIjowLjAwMDY1NjE5OTM1MDQ4NzQ0MDh9XX0=`

const (
	Greet   = "greet"
	GoodBye = "goodbye"
	Affirm  = "affirm"
	Deny    = "deny"
)

var numberRestMap *PhonenumberMap

func main() {

	e := echo.New()

	numberRestMap = new(PhonenumberMap)

	e.GET("/v1/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "Healthy!!!")
	})

	e.POST("/SiprtcApplications/VoxvellyDemo", VoxvellyDemo)
	e.POST("/SiprtcApplications/VoxvellyDemoDtmfReceived", VoxvellyDemoDtmfReceived)

	e.POST("/SiprtcApplications/DirectCall", DirectCall)

	e.GET("/SiprtcApplications/MainRestaurantMenu", MainRestaurantMenu)
	e.POST("/SiprtcApplications/DtmfReceived", RestaurantDtmfReceived)

	e.POST("/SiprtcApplications/VoicebotLoan", func(c echo.Context) error {
		u := StatusCallback{}
		err := c.Bind(&u)

		rejectresp := GetRejectedResponse()

		if err != nil {
			return c.XML(http.StatusOK, rejectresp)
		}

		ivrRest := numberRestMap.GetNumberInstance(u.From)

		if ivrRest == nil {
			rejectresp := GetRejectedResponse()
			return c.XML(http.StatusOK, rejectresp)
		}

		return c.XML(http.StatusOK, ivrRest.CreateWelcomeVoiceBot("Welcome to ICICI Bank, Are you interested in home loan? Say YES or NO"))
	})

	e.POST("/SiprtcApplications/Voicebot", func(c echo.Context) error {
		u := StatusCallback{}
		err := c.Bind(&u)

		rejectresp := GetRejectedResponse()

		if err != nil {
			return c.XML(http.StatusOK, rejectresp)
		}

		fmt.Println("Voicebot FROM : ", u.From)

		ivrRest := new(RestaurentIVR)
		numberRestMap.StoreNumberInstance(u.From, ivrRest)

		ssmlText := `<speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" xmlns:emo="http://www.w3.org/2009/10/emotionml" version="1.0" xml:lang="en-US"><voice name="en-IN-NeerjaNeural"><prosody rate="0%" pitch="0%">Welcome to Big Pitcher, How can i help you?</prosody></voice></speak>`
		return c.XML(http.StatusOK, ivrRest.CreateWelcomeVoiceBot(ssmlText))
	})

	e.POST("/SiprtcApplications/UserIntent", func(c echo.Context) error {
		resp := &Response{}
		resp.Text = ""
		u := StatusCallback{}
		err := c.Bind(&u)
		if err != nil {
			rejectresp := &Response{}
			rejectresp.Reject = &Reject{
				Reason: "rejected",
			}
			return c.XML(http.StatusOK, resp)
		}

		userIntent := ProcessUserIntent(u.UserIntent)

		fmt.Println("Voicebot FROM : ", u.From)

		ivrRest := numberRestMap.GetNumberInstance(u.From)

		if ivrRest == nil {
			rejectresp := GetRejectedResponse()
			return c.XML(http.StatusOK, rejectresp)
		}

		fmt.Println(userIntent)

		prefix := `<speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" xmlns:emo="http://www.w3.org/2009/10/emotionml" version="1.0" xml:lang="en-US"><voice name="en-IN-NeerjaNeural"><prosody rate="0%" pitch="0%">`
		postfix := `</prosody></voice></speak>`
		switch userIntent.Intent.Name {
		case "greet":
			fmt.Println("greet")
			ssmlText := prefix + `Hello` + postfix
			resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
		case "goodbye":
			fmt.Println("goodbye")
			ssmlText := prefix + `Good Bye, it was nice talking to you.` + postfix
			resp = CreateSayHangupSSML(ssmlText)
		case "affirm":
			fmt.Println("affirm")
			ssmlText := prefix + `Thank you! for confirming that you are interested in booking table. how can i help you with booking?` + postfix
			resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
		case "deny":
			fmt.Println("deny")
			ssmlText := prefix + `Ok No problem, we will not call you back again. Thanks for your feedback on product. Bye` + postfix
			resp = CreateSayHangupSSML(ssmlText)
		case "complain":
			fmt.Println("complain")
			ssmlText := prefix + `Let me help you in raising complain. Please speak out for next 5 minute, your audio is recorded and on priority your complaing would be analysed and resolved.` + postfix
			resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
		case "booking":
			fmt.Println("booking")
			ssmlText := prefix + `I can help you with booking of table, For how many persons do you need reservation?` + postfix
			resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
		case "booking_with_count":
			fmt.Println("booking_with_count")
			if userIntent.Entities[0].Entity == "number" {
				fmt.Println("number of persongs : ", userIntent.Entities[0].Value)
				// ask for what time today or tomorrow?
				ivrRest.SetCount(userIntent.Entities[0].Value)
				ssmlText := prefix + `I would like to confirm that you need booking for ` + userIntent.Entities[0].Text + `, and do you need booking today, tomorrow or day after tomorrow?` + postfix
				resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
			}
		case "booking_with_count_time":
			fmt.Println("booking_with_count_time")
			// ask for the
		case "booking_time_day":
			fmt.Println("booking_time_day")
			if userIntent.Entities[0].Entity == "time" {
				fmt.Println("time of booking : ", userIntent.Entities[0].Value)
				ivrRest.SetDayTime(userIntent.Entities[0].Text)
				ssmlText := prefix + `Ok i will do booking ` + userIntent.Entities[0].Text + ` for you, What time do you need booking? you can say like 9:30PM or 10:30AM.` + postfix
				resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
			}
			// ivrRest.CreateWelcomeVoiceBot(ssmlText)
		case "booking_with_time_day_hours_minute":
			fmt.Println("booking_with_time_day_hours_minute")
			if userIntent.Entities[0].Entity == "time" {
				fmt.Println("time of booking : ", userIntent.Entities[0].Text, userIntent.Entities[0].Value)
				// ask for the time of today or tomorrow?
				ssmlText := prefix + `Ok i will do booking ` + userIntent.Entities[0].Text + ` for you.` + postfix
				resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
			}
		case "booking_with_count_time_day_hours_minute":
			fmt.Println("booking_with_count_time_day_hours_minute")
			for i := 0; i < len(userIntent.Entities); i++ {
				if userIntent.Entities[0].Entity == "time" {
					fmt.Println("time of booking : ", userIntent.Entities[0].Text, userIntent.Entities[0].Value)
					ivrRest.SetDayTime(userIntent.Entities[0].Text)
				} else if userIntent.Entities[0].Entity == "number" {
					fmt.Println("number of persongs : ", userIntent.Entities[0].Value)
				}
			}
			ssmlText := prefix + `Ok i will do booking ` + userIntent.Entities[0].Text + ` for you.` + postfix
			resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)

		case "booking_count":
			fmt.Println("booking_count")
			if userIntent.Entities[0].Entity == "number" {
				fmt.Println("number of persongs : ", userIntent.Entities[0].Value)
				ivrRest.SetCount(userIntent.Entities[0].Value)
				// ask for what time today or tomorrow?
				ssmlText := prefix + `I would like to confirm that you need booking for ` + userIntent.Entities[0].Text + `, and do you need booking today, tomorrow or day after tomorrow?` + postfix
				resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
			}
		case "talkto_agent":
			fmt.Println("talkto_agent")
			ssmlText := prefix + `Ok , Let me transfer your call to an agent. Transferring call now.` + postfix
			resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
		case "nlu_fallback":
			ssmlText := prefix + `I don't understand can you please say that again?` + postfix
			resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
		default:
			fmt.Println("Invalid")
		}

		return c.XML(http.StatusOK, resp)
	})
	e.Logger.Fatal(e.Start(":8080"))
}
