package main

import (
	"fmt"

	"github.com/notedit/sdp"
	"github.com/notedit/webrtc-transport"
)

var Capabilities = map[string]*sdp.Capability{
	"audio": &sdp.Capability{
		Codecs: []string{"opus"},
	},
	"video": &sdp.Capability{
		Codecs: []string{"h264"},
		Rtx:    true,
		Rtcpfbs: []*sdp.RtcpFeedback{
			&sdp.RtcpFeedback{
				ID: "transport-cc",
			},
			&sdp.RtcpFeedback{
				ID:     "ccm",
				Params: []string{"fir"},
			},
			&sdp.RtcpFeedback{
				ID:     "nack",
				Params: []string{"pli"},
			},
		},
		Extensions: []string{
			"urn:3gpp:video-orientation",
			"http://www.ietf.org/id/draft-holmer-rmcat-transport-wide-cc-extensions-01",
			"http://www.webrtc.org/experiments/rtp-hdrext/abs-send-time",
			"urn:ietf:params:rtp-hdrext:toffse",
			"urn:ietf:params:rtp-hdrext:sdes:rtp-stream-id",
			"urn:ietf:params:rtp-hdrext:sdes:mid",
		},
	},
}

func main() {

	fmt.Println(Capabilities)

	webrtcTransport := webrtc.NewWebRTCTransport()

	fmt.Println(webrtcTransport)
}
