package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func DirectCall(c echo.Context) error {
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
	resp.Dial = &Dial{
		AnswerOnBridge: true,
		Number: &Number{
			Text: u.To,
		},
	}
	return c.XML(http.StatusOK, resp)
}

/*
Welcome to the voxvalley technologies Press 1 for sales Press 2 for support

1. forward call to suren
2. forward call to naresh.
*/
func DtmfDemo(c echo.Context) error {
	actionURL := "https://demo.siprtc.io/SiprtcApplications/DtmfDemoDtmfReceived"
	resp := CreateGatherSayResponse("Welcome to the voxvalley technologies. Press 1 for sales, Press 2 for support.", actionURL, "1")
	return c.XML(http.StatusOK, resp)
}

func DtmfDemoDtmfReceived(c echo.Context) error {
	var resp *Response

	u := StatusCallback{}
	err := c.Bind(&u)
	rejectresp := GetRejectedResponse()

	if err != nil {
		return c.XML(http.StatusOK, rejectresp)
	}

	fmt.Println("Dtmf Digit : ", u.Digits, u.From)

	if len(u.Digits) == 0 {
		return c.XML(http.StatusOK, rejectresp)
	}

	lastChar := u.Digits[len(u.Digits)-1:]

	if lastChar == "#" {
		u.Digits = u.Digits[0 : len(u.Digits)-1]
	}

	if u.Digits == "1" {
		resp = CreateSayDial("Forwarding call to sales", "91XXXXXXXXXX", "+91XXXXXXXXXX")
	} else {
		resp = CreateSayDial("Forwarding call to support", "91XXXXXXXXXX", "+91XXXXXXXXXX")
	}

	return c.XML(http.StatusOK, resp)
}

func MainRestaurantMenu(c echo.Context) error {
	u := StatusCallback{}
	err := c.Bind(&u)
	if err != nil {
		rejectresp := &Response{}
		rejectresp.Reject = &Reject{
			Reason: "rejected",
		}
		return c.XML(http.StatusOK, rejectresp)
	}
	ivrRest := new(RestaurentIVR)
	ivrRest.InitRestaurentIVR()
	numberRestMap.StoreNumberInstance(u.From, ivrRest)
	resp := ivrRest.GetMainMenuResponse()
	return c.XML(http.StatusOK, resp)
}

func RestaurantDtmfReceived(c echo.Context) error {
	u := StatusCallback{}
	err := c.Bind(&u)
	rejectresp := GetRejectedResponse()

	if err != nil {
		return c.XML(http.StatusOK, rejectresp)
	}

	fmt.Println("Dtmf Digit : ", u.Digits, u.From)

	if len(u.Digits) == 0 {
		return c.XML(http.StatusOK, rejectresp)
	}

	lastChar := u.Digits[len(u.Digits)-1:]

	if lastChar == "#" {
		u.Digits = u.Digits[0 : len(u.Digits)-1]
	}

	ivrRest := numberRestMap.GetNumberInstance(u.From)

	if ivrRest == nil {
		rejectresp := GetRejectedResponse()
		return c.XML(http.StatusOK, rejectresp)
	}

	resp := ivrRest.ProcessDTMFDigits(u.Digits)

	return c.XML(http.StatusOK, resp)
}
