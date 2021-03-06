package webrtc

import (
	"fmt"
	"github.com/notedit/sdp"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/internal/mux"
)

var api = webrtc.NewAPI()

type WebRTCTansport struct {
	api           *webrtc.API
	gatherer      *webrtc.ICEGatherer
	iceTransport  *webrtc.ICETransport
	dtlsTransport *webrtc.DTLSTransport

	localCandidates []*sdp.CandidateInfo
}

func NewWebRTCTransport() *WebRTCTansport {

	iceOptions := webrtc.ICEGatherOptions{
		ICEServers: []webrtc.ICEServer{},
	}

	gatherer, err := api.NewICEGatherer(iceOptions)
	if err != nil {
		panic(err)
	}

	fmt.Println(mux.MatchAll)

	ice := api.NewICETransport(gatherer)

	// Construct the DTLS transport
	dtls, err := api.NewDTLSTransport(ice, nil)
	if err != nil {
		panic(err)
	}

	webrtcTransport := &WebRTCTansport{
		api:             api,
		iceTransport:    ice,
		dtlsTransport:   dtls,
		gatherer:        gatherer,
		localCandidates: make([]*sdp.CandidateInfo, 0),
	}
	return webrtcTransport
}

func (t *WebRTCTansport) GetLocalICEInfo() (*sdp.ICEInfo, error) {
	iceParams, err := t.gatherer.GetLocalParameters()
	if err != nil {
		return nil, err
	}
	iceInfo := sdp.NewICEInfo(iceParams.UsernameFragment, iceParams.Password)
	iceInfo.SetLite(true)
	return iceInfo, nil
}

func (t *WebRTCTansport) GetLocalDTLSInfo() (*sdp.DTLSInfo, error) {
	dtlsParams, err := t.dtlsTransport.GetLocalParameters()
	if err != nil {
		return nil, err
	}

	fmt.Println(dtlsParams.Fingerprints)

	fingerprint := dtlsParams.Fingerprints[0]

	fmt.Println(fingerprint)

	var setup sdp.Setup
	if dtlsParams.Role == webrtc.DTLSRoleClient {
		setup = sdp.SETUPACTIVE
	} else if dtlsParams.Role == webrtc.DTLSRoleServer {
		setup = sdp.SETUPPASSIVE
	} else if dtlsParams.Role == webrtc.DTLSRoleAuto {
		setup = sdp.SETUPACTPASS
	}

	dtlsInfo := sdp.NewDTLSInfo(setup, fingerprint.Algorithm, fingerprint.Value)
	return dtlsInfo, nil
}

func (t *WebRTCTansport) GetLocalCandidates() ([]*sdp.CandidateInfo, error) {

	err := t.gatherer.Gather()
	if err != nil {
		return nil, err
	}

	candidates, err := t.gatherer.GetLocalCandidates()
	if err != nil {
		return nil, err
	}

	for _, candidate := range candidates {
		candidateInfo := sdp.NewCandidateInfo(candidate.Foundation, int(candidate.Component), candidate.Protocol.String(), int(candidate.Priority), candidate.Address, int(candidate.Port), candidate.Typ.String(), "", 0)
		t.localCandidates = append(t.localCandidates, candidateInfo)
	}

	fmt.Println(t.localCandidates)

	return t.localCandidates, nil
}

func (t *WebRTCTansport) SetRemoteICEInfo(ice *sdp.ICEInfo) error {

	iceParams := webrtc.ICEParameters{
		UsernameFragment: ice.GetUfrag(),
		Password:         ice.GetPassword(),
		ICELite:          false,
	}
	iceRole := webrtc.ICERoleControlling
	err := t.iceTransport.Start(nil, iceParams, &iceRole)
	return err
}

func (t *WebRTCTansport) AddRemoteCandidate(candidate *sdp.CandidateInfo) {

	protocol, _ := webrtc.NewICEProtocol(candidate.GetTransport())
	typ, _ := webrtc.NewICECandidateType(candidate.GetType())

	can := webrtc.ICECandidate{
		Foundation: candidate.GetFoundation(),
		Priority:   uint32(candidate.GetPriority()),
		Address:    candidate.GetAddress(),
		Protocol:   protocol,
		Port:       uint16(candidate.GetPort()),
		Component:  uint16(candidate.GetComponentID()),
		Typ:        typ,
	}

	t.iceTransport.AddRemoteCandidate(can)

}

func (t *WebRTCTansport) SetRemoteDTLSInfo(dtls *sdp.DTLSInfo) error {

	var role webrtc.DTLSRole
	if dtls.GetSetup() == sdp.SETUPACTIVE {
		role = webrtc.DTLSRoleClient
	} else if dtls.GetSetup() == sdp.SETUPPASSIVE {
		role = webrtc.DTLSRoleServer
	} else {
		role = webrtc.DTLSRoleAuto
	}

	fingerprint := webrtc.DTLSFingerprint{
		Algorithm: dtls.GetHash(),
		Value:     dtls.GetFingerprint(),
	}

	dtlsParams := webrtc.DTLSParameters{
		Role:         role,
		Fingerprints: []webrtc.DTLSFingerprint{fingerprint},
	}

	err := t.dtlsTransport.Start(dtlsParams)
	return err
}
