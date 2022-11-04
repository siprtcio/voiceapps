package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/labstack/echo"
)

type UserIntent struct {
	Text          string          `json:"text"`
	Intent        Intent          `json:"intent"`
	Entities      []Entities      `json:"entities"`
	IntentRanking []IntentRanking `json:"intent_ranking"`
}
type Intent struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}
type Entities struct {
	Start      int             `json:"start"`
	End        int             `json:"end"`
	Text       string          `json:"text"`
	Value      json.RawMessage `json:"value"`
	Confidence int             `json:"confidence"`
	Entity     string          `json:"entity"`
}
type IntentRanking struct {
	Name       string  `json:"name"`
	Confidence float64 `json:"confidence"`
}

type ValueType int64

const (
	VAL_INVALID ValueType = 0
	VAL_INT               = 1
	VAL_STRING            = 2
	VAL_STRUCT            = 3
)

type ValueStruct struct {
	To   time.Time `json:"to"`
	From time.Time `json:"from"`
}

func stringOrIntOrStruct(Value []byte) (int, string, *ValueStruct, ValueType) {
	if utf8.Valid(Value) {
		i, err := strconv.Atoi(string(Value))
		if err == nil {
			fmt.Println("got int: " + strconv.Itoa(i))
			return i, "", nil, VAL_INT
		} else {
			valStruct := ValueStruct{}
			err := json.Unmarshal(Value, &valStruct)
			if err != nil {
				fmt.Println("got string", string(Value))
				return 0, string(Value), nil, VAL_STRING
			} else {
				fmt.Println("got struct", valStruct)
				return 0, "", &valStruct, VAL_STRUCT
			}
		}
	} else {
		fmt.Println("whoops")
	}
	return 0, "", nil, VAL_INVALID
}

func ProcessUserIntent(uIntent string) UserIntent {
	// userIntent is base64 encoded data.
	data, err := base64.StdEncoding.DecodeString(uIntent)
	if err != nil {
		log.Fatal("error:", err)
	}
	fmt.Println(string(data))
	// decode it to json
	var userIntent UserIntent
	err = json.Unmarshal(data, &userIntent)
	if err != nil {
		fmt.Println(err)
	}
	// parse json and get the intent
	fmt.Println(userIntent.Text, userIntent.Intent.Name)
	// process intent.
	return userIntent
}

func VoicebotWelcomeMessage(c echo.Context) error {
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
}

func VoicebotUserIntent(c echo.Context) error {
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

	textlength := len(userIntent.Text)

	if textlength == 0 || userIntent.Intent.Confidence < 0.85 {
		userIntent.Intent.Name = "nlu_fallback"
	}

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
	case "cancel_booking":
		fmt.Println("complain")
		ssmlText := prefix + `Let me help you in cancelling reservation. As i can see your current reservation is there. I am cancelling now and you will get sms for the same. Thanks for your time..` + postfix
		resp = CreateSayHangupSSML(ssmlText)
	case "booking":
		fmt.Println("booking")
		ssmlText := prefix + `I can help you with booking of table, For how many persons do you need reservation?` + postfix
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
	case "booking_with_count":
		fmt.Println("booking_with_count")
		ssmlText := ""
		if len(userIntent.Entities) == 0 {
			ssmlText = prefix + `I'm sorry, But looks like you are trying to say number of person visiting us, I don't understand, can you please say that again?` + postfix
		} else {
			if userIntent.Entities[0].Entity == "number" {
				fmt.Println("number of persongs : ", userIntent.Entities[0].Value)
				// ask for what time today or tomorrow?
				count, _, _, _ := stringOrIntOrStruct(userIntent.Entities[0].Value)
				ivrRest.SetCount(count)
				ssmlText = prefix + `I would like to confirm that you need booking for ` + userIntent.Entities[0].Text + `, and do you need booking today, tomorrow or day after tomorrow?` + postfix
			}
		}
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)

	case "booking_with_count_time":
		fmt.Println("booking_with_count_time")
		ssmlText := ""
		if len(userIntent.Entities) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand, can you please say that again?` + postfix
		} else {
			for i := 0; i < len(userIntent.Entities); i++ {
				if userIntent.Entities[0].Entity == "time" {
					fmt.Println("time of booking : ", userIntent.Entities[0].Text, userIntent.Entities[0].Value)
					ivrRest.SetDayTime(userIntent.Entities[0].Text)
					ssmlText = prefix + `I would like to confirm that you need booking ` + userIntent.Entities[0].Text
				} else if userIntent.Entities[0].Entity == "number" {
					count, _, _, _ := stringOrIntOrStruct(userIntent.Entities[0].Value)
					fmt.Println("number of persongs : ", count)
					ivrRest.SetCount(count)
					ssmlText = ssmlText + `for ` + strconv.Itoa(count) + `people`
				}
			}
			ssmlText = ssmlText + `What time do you need booking? you can say like 9:30PM or 10:30AM.` + postfix
		}
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
		// ask for the
	case "booking_time_day":
		fmt.Println("booking_time_day")
		ssmlText := ""
		if len(userIntent.Entities) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand time correctly, can you please say that again?` + postfix
		} else {
			if userIntent.Entities[0].Entity == "time" {
				fmt.Println("time of booking : ", userIntent.Entities[0].Value)
				ivrRest.SetDayTime(userIntent.Entities[0].Text)
				ssmlText = prefix + `Ok i will do booking ` + userIntent.Entities[0].Text + ` for you, What time do you need booking? you can say like 9:30PM or 10:30AM.` + postfix
			}
		}

		if len(ssmlText) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand, can you please say that again?` + postfix
		}
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
	case "booking_with_time_day_hours_minute":
		fmt.Println("booking_with_time_day_hours_minute")
		ssmlText := ""
		if len(userIntent.Entities) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand time correctly, can you please say that again?` + postfix
		} else {
			if userIntent.Entities[0].Entity == "time" {
				fmt.Println("time of booking : ", userIntent.Entities[0].Text, userIntent.Entities[0].Value)
				// ask for the time of today or tomorrow?
				ssmlText = prefix + `Ok i will do booking ` + userIntent.Entities[0].Text + ` for you.` + postfix
				ivrRest.SetDayTime(userIntent.Entities[0].Text)
			}
		}
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)

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
		ssmlText := ""
		if len(userIntent.Entities) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand, can you please say that again?` + postfix
		} else {
			if userIntent.Entities[0].Entity == "number" {
				count, _, _, _ := stringOrIntOrStruct(userIntent.Entities[0].Value)
				fmt.Println("number of persongs : ", count)
				ivrRest.SetCount(count)
				// ask for what time today or tomorrow?
				ssmlText = prefix + `I would like to confirm that you need booking for ` + userIntent.Entities[0].Text + `, and do you need booking today, tomorrow or day after tomorrow?` + postfix
			} else {
				ssmlText = prefix + `I'm sorry, I don't understand, can you please say that again?` + postfix
			}
		}

		if len(ssmlText) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand, can you please say that again?` + postfix
		}
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)

	case "booking_time":
		fmt.Println("booking_time")
		ssmlText := ""
		if len(userIntent.Entities) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand, can you please say that again?` + postfix
		} else {
			if userIntent.Entities[0].Entity == "time" {
				_, bookingTime, bookingTimeStruct, valType := stringOrIntOrStruct(userIntent.Entities[0].Value)
				if valType == VAL_STRING {
					ivrRest.SetDayTime(bookingTime)
				} else {
					ivrRest.SetBookingDayTime(bookingTimeStruct.From)
				}
				ssmlText = prefix + `I confirm your booking ` + userIntent.Entities[0].Text + `, and do you will get sms for confirmation of booking.` + postfix
			}
		}

		if len(ssmlText) == 0 {
			ssmlText = prefix + `I'm sorry, I don't understand, can you please say that again?` + postfix
		}
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)

	case "talkto_agent":
		fmt.Println("talkto_agent")
		ssmlText := prefix + `Ok , Let me transfer your call to an agent. Transferring call now.` + postfix
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
	case "nlu_fallback":
		ssmlText := ""
		if textlength == 0 {
			ssmlText = prefix + `I didn't hear you, can you please say that again?` + postfix
		} else {
			ssmlText = prefix + `I don't understand can you please say that again?` + postfix
		}
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
	default:
		ssmlText := prefix + `I'm sorry, I didn't catch that. Would you mind speaking more slowly?` + postfix
		resp = ivrRest.CreateWelcomeVoiceBot(ssmlText)
	}
	return c.XML(http.StatusOK, resp)
}
