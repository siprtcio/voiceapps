package main

func GetRejectedResponse() *Response {
	rejectresp := &Response{}
	rejectresp.Reject = &Reject{
		Reason: "rejected",
	}
	return rejectresp
}

func CreateSayDial(sayText string, dialNumber string) *Response {
	resp := &Response{}
	resp.Text = ""
	resp.Say = &Say{
		Text: sayText,
	}

	resp.Dial = &Dial{}
	resp.Dial.Text = dialNumber
	return resp
}

func CreateGatherSayResponse(gatherSayString string, actionURL string, digits string) *Response {
	resp := &Response{}
	resp.Text = ""
	resp.Gather = &Gather{
		Action:      actionURL,
		NumDigits:   digits,
		FinishOnKey: "#",
		Method:      "POST",
	}

	resp.Gather.Say = &Say{
		Text: gatherSayString,
	}

	resp.Say = &Say{
		Text: "We didn't receive any input. Goodbye!",
	}
	return resp
}

func CreateSayHangup(sayString string) *Response {
	resp := &Response{}
	resp.Text = ""
	resp.Say = &Say{
		Text: sayString,
	}
	resp.Hangup = &Hangup{}
	return resp
}
