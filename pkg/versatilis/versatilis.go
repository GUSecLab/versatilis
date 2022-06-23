package versatilis

import (
	_ "embed"

	"github.com/Masterminds/semver"
)

//go:embed version.md
var versionString string

var Version *semver.Version

type VersatilisErrNotImplemented struct{}

func (e *VersatilisErrNotImplemented) Error() string {
	return "not yet implemented"
}

/*
type State struct {
	Name                     string
	initiator                bool
	noiseConfig              *noise.Config
	noiseHandshakeState      *noise.HandshakeState
	handShakeCompleted       bool
	handshakeSendRecvPattern []bool // "true"s indicate send events; "false" are receives
	encryptState             *noise.CipherState
	decryptState             *noise.CipherState
	connected                bool
	dst                      *Address
	listenAddress            *Address
	//readBuffer               readBuffer
	//readMu *sync.Mutex
}
*/

func init() {
	var err error
	Version, err = semver.NewVersion(versionString)
	if err != nil {
		panic(err)
	}
}

/*
func New(initiator bool, name string) *State {
	var err error

	state := new(State)
	state.Name = name
	state.initiator = initiator
	state.handShakeCompleted = false
	state.connected = false

	state.noiseConfig = &noise.Config{
		CipherSuite: noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256),
		Random:      rand.Reader,
		Initiator:   initiator,
		Pattern:     noise.HandshakeNN, // TODO
	}

	if state.initiator { // TODO: this is specific to noise.HandshakeNN
		z := [...]bool{true, false}
		state.handshakeSendRecvPattern = z[:]
	} else {
		z := [...]bool{false, true}
		state.handshakeSendRecvPattern = z[:]
	}

	if state.noiseHandshakeState, err = noise.NewHandshakeState(*state.noiseConfig); err != nil {
		panic(err)
	}
	log.Debugf("[%v] current handshake state message is %v", state.Name, state.noiseHandshakeState.MessageIndex())

	return state
}

func SetLogLevel(level log.Level) {
	log.SetLevel(level)
}
*/

/*
func (state *State) DoHandshake(dst *Address, listenAddress *Address) {
	var out2 []byte
	var err error

	out := make([]byte, 0, 4096)

	for _, action := range state.handshakeSendRecvPattern {
		if action { // send event
			out2, state.encryptState, state.decryptState, err =
				state.noiseHandshakeState.WriteMessage(out, nil)
			log.Debugf("[%v] sending message of length %v", state.Name, len(out2))

			if err := sendPackageViaTransport(dst, &Package{
				Version:            Version.String(),
				NoiseHandshakeInfo: out2,
			}); err != nil {
				panic(err)
			}
		} else { // receive event
			log.Debugf("[%v] waiting to receive message", state.Name)
			var p *Package
			p, err = recvPackageViaTransport(listenAddress, true)
			if err != nil {
				panic(err)
			}
			log.Debugf("[%v] message received of length %v", state.Name, len(p.NoiseHandshakeInfo))
			_, state.decryptState, state.encryptState, err =
				state.noiseHandshakeState.ReadMessage(out, p.NoiseHandshakeInfo)
		}
	}

	if err != nil {
		panic(err)
	}
	if state.encryptState != nil {
		state.handShakeCompleted = true
		state.dst = dst
		state.listenAddress = listenAddress
		log.Infof("[%v] handshake complete", state.Name)
		log.Debugf("[%v] enc_state=%v; dec_state=%v", state.Name, state.encryptState, state.decryptState)
	}
}
*/
