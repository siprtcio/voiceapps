package main

import (
	"bytes"
	"encoding/xml"
	"net/http"

	"github.com/labstack/echo"
)

var samplexml = []byte(`<?xml version="1.0"?>
<Application>
<Node nodeId="1" nodeType="gather" >
<Condition></Condition>
<Response><Gather action="https://demo.siprtc.io/SiprtcApplications/DtmfDemoDtmfReceived" method="POST" finishOnKey="#" numDigits="1"><Say>Welcome to the voxvalley technologies. Press 1 for sales, Press 2 for support.</Say></Gather><Say>We didn&#39;t receive any input. Goodbye!</Say></Response>
<ChildrenNodes>
<Node nodeId="2" nodeType="sayDial">
<Condition condtionType="dtmf">
1
</Condition>
<Response>
<Say>Forwarding your call to support</Say>
<Dial>9945073606</Dial>
</Response>
<ChildrenNodes></ChildrenNodes>
</Node>
<Node nodeId="3" nodeType="sayDial">
<Condition condtionType="dtmf">
2
</Condition>
<Response>
<Say>Forwarding your call to sales</Say>
<Dial>9945073606</Dial>
</Response>
<ChildrenNodes></ChildrenNodes>
</Node>
</ChildrenNodes>
</Node>
</Application>
`)

type ApplicationContext struct {
	currentNode *Node
	variables   [][]string
}

var currentNode *Node

type Condition struct {
	XMLName       xml.Name `xml:"Condition"`
	ConditionType string   `xml:"conditionType,attr,omitempty"`
	Text          string   `xml:",chardata"`
}

type ChildrenNodes struct {
	XMLName       xml.Name `xml:"ChildrenNodes"`
	Text          string   `xml:",chardata"`
	ChildrenNodes []*Node  `xml:"ChildrenNodes,omitempty"`
}

type Node struct {
	XMLName       xml.Name       `xml:"Node"`
	Text          string         `xml:",chardata"`
	NodeID        string         `xml:"nodeId,attr,omitempty"`
	NodeType      string         `xml:"nodeType,attr,omitempty"`
	Condition     *Condition     `xml:"Condition,omitempty"`
	ChildrenNodes *ChildrenNodes `xml:"ChildrenNodes,omitempty"`
	Response      *Response      `xml:"Response"`
}

type Application struct {
	XMLName xml.Name `xml:"Application"`
	Text    string   `xml:",chardata"`
	Node    *Node    `xml:"Node,omitempty"`
}

func ProcessApplicationWelcomeXML(c echo.Context) error {
	buf := bytes.NewBuffer(samplexml)
	dec := xml.NewDecoder(buf)
	var app Application
	err := dec.Decode(&app)
	if err != nil {
		panic(err)
	}
	currentNode = app.Node
	return c.XML(http.StatusOK, app.Node.Response)
}

func ProcessApplicationDtmf(c echo.Context) error {
	if currentNode.NodeType == "gather" {

	}
	return c.XML(http.StatusOK, nil)
}

func ProcessApplicationIntent(c echo.Context) error {
	if currentNode.NodeType == "gather" {

	}
	return c.XML(http.StatusOK, nil)
}
