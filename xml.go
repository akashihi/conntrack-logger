package main

import (
    "encoding/xml"
    "log"
    "time"
)

  /*
    <flow type="new">
	<meta direction="original">
	    <layer3 protonum="2" protoname="ipv4">
		<src>192.168.58.3</src>
		<dst>192.168.58.4</dst>
	    </layer3>
	    <layer4 protonum="6" protoname="tcp">
		<sport>51156</sport>
		<dport>61616</dport>
	    </layer4>
	</meta>
	<meta direction="reply">
	    <layer3 protonum="2" protoname="ipv4">
		<src>192.168.58.4</src>
		<dst>192.168.58.3</dst>
	    </layer3>
	    <layer4 protonum="6" protoname="tcp">
		<sport>61616</sport>
		<dport>51156</dport>
	    </layer4>
	</meta>
	<meta direction="independent">
	    <state>SYN_SENT</state>
	    <timeout>120</timeout>
	    <id>1764685040</id>
	    <unreplied/>
	</meta>
    </flow>
  */

type NFL3 struct {
    Protonum  int        `xml:"protonum,attr"`
    Protoname string     `xml:"protoname,attr"`
    Src       string	 `xml:"src"`
    Dst       string     `xml:"dst"`
}

type NFL4 struct {
    Protonum  int        `xml:"protonum,attr"`
    Protoname string     `xml:"protoname,attr"`
    Sport     int        `xml:"sport"`
    Dport     int        `xml:"dport"`
}

type NMeta struct {
    Direction string     `xml:"direction,attr"`
    Layer3    NFL3	 `xml:"layer3"`
    Layer4    NFL4       `xml:"layer4"`
}

type NFlow struct {
    Type      string     `xml:"type,attr"`
    Metas     []NMeta    `xml:"meta"`
}

func parseNF(messages <-chan []byte, flows chan<- FlowRecord) {
    for {
        m := <-messages

	var flow NFlow
        err := xml.Unmarshal(m, &flow)
	if err != nil {
	    log.Fatal(err)
            continue
        }

        for _, meta := range flow.Metas {
    	    if meta.Direction == "original" {
    		result := FlowRecord { 
    		    TS : time.Now(),
                    Proto : meta.Layer4.Protoname,
            	    Src : meta.Layer3.Src,
            	    Dst : meta.Layer3.Dst,
            	    Sport : meta.Layer4.Sport,
            	    Dport : meta.Layer4.Dport,
		}
        	flows <- result
        	break
    	    }
	}
    }
}