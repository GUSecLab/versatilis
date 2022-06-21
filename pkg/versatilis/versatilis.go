package versatilis

import (
	"crypto/rand"
	_ "embed"
	"encoding/json"

	"github.com/Masterminds/semver"
	"github.com/flynn/noise"
	log "github.com/sirupsen/logrus"
)

//go:embed version.md
var versionString string

var Version *semver.Version

type KeyType int64

const (
	KeyTypeStatic KeyType = iota
	KeyTypeEphemeral
)

func (k KeyType) String() string {
	switch k {
	case KeyTypeEphemeral:
		return "ephemeral"
	case KeyTypeStatic:
		return "static"
	default:
		panic("no such key type")
	}
}

type State struct {
	Name                     string
	initiator                bool
	noiseConfig              *noise.Config
	noiseHandshakeState      *noise.HandshakeState
	handShakeCompleted       bool
	handshakeSendRecvPattern []bool // "true"s indicate send events; "false" are receives
	encryptState             *noise.CipherState
	decryptState             *noise.CipherState
}

type Message struct {
	Id      string      `json:"id"`
	Payload interface{} `json:"payload"`
}

type MessageBuffer []*Message

func init() {
	var err error
	Version, err = semver.NewVersion(versionString)
	if err != nil {
		panic(err)
	}
}

func New(initiator bool, name string) *State {
	var err error

	state := new(State)
	state.Name = name
	state.initiator = initiator
	state.handShakeCompleted = false

	state.noiseConfig = &noise.Config{
		CipherSuite: noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256),
		Random:      rand.Reader,
		Initiator:   initiator,
		Pattern:     noise.HandshakeNN, // TODO
	}

	if state.initiator {
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

func (state *State) DoHandshake(dst *Address, listenAddress *Address) {
	//	sendChannel chan *Package, recvChannel chan *Package) {
	var out2 []byte
	var err error

	out := make([]byte, 0, 4096)

	for _, action := range state.handshakeSendRecvPattern {
		if action { // send event
			out2, state.encryptState, state.decryptState, err =
				state.noiseHandshakeState.WriteMessage(out, nil)
			log.Debugf("[%v] sending message of length %v", state.Name, len(out2))

			if err := send(dst, &Package{
				Version:            Version.String(),
				NoiseHandshakeInfo: out2,
			}); err != nil {
				panic(err)
			}
		} else { // receive event
			log.Debugf("[%v] waiting to receive message", state.Name)
			var p *Package
			p, err = recv(listenAddress, true)
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
		log.Infof("[%v] handshake complete", state.Name)
		log.Debugf("[%v] enc_state=%v; dec_state=%v", state.Name, state.encryptState, state.decryptState)
	}
}

func (state *State) Send(dst *Address, buffer *MessageBuffer) error {
	for _, message := range *buffer {
		plaintext, err := json.Marshal(message)
		if err != nil {
			return err
		}
		//ciphertext := make([]byte, 0)
		ad := make([]byte, 16) // TODO
		rand.Read(ad)

		ciphertext, err := state.encryptState.Encrypt(nil, ad, plaintext)
		if err != nil {
			return err
		}
		log.Debugf("sending message %v", message)

		if err := send(dst, &Package{
			Version:         Version.String(),
			NoiseCiphertext: ciphertext,
			NoiseAuthTag:    ad,
		}); err != nil {
			return err
		}
	}
	return nil
}

// receives a message or returns nil, nil if no message is available (if block is false)
func (state *State) Receive(listenAddress *Address, block bool) (*Message, error) {

	p, err := recv(listenAddress, block)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, nil
	}

	plaintext, err := state.decryptState.Decrypt(nil, p.NoiseAuthTag, p.NoiseCiphertext)
	if err != nil {
		return nil, err
	}
	var message Message
	err = json.Unmarshal(plaintext, &message)
	if err != nil {
		return nil, err
	} else {
		return &message, nil
	}
}
