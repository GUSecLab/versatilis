package versatilis

import (
	"crypto/rand"

	"github.com/flynn/noise"
	log "github.com/sirupsen/logrus"
)

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
	name        string
	noiseConfig *noise.Config
}

type Message struct {
	Id      string
	Payload interface{}
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.Debug("initializing Versātilis")
}

func New(initiator bool, name string) *State {
	state := new(State)
	state.name = name

	state.noiseConfig = &noise.Config{
		CipherSuite: noise.NewCipherSuite(noise.DH25519, noise.CipherAESGCM, noise.HashSHA256),
		Random:      rand.Reader,
		Initiator:   initiator,
		Pattern:     noise.HandshakeNN, // TODO
	}
	return state
}

func (state *State) GenKey(keyType KeyType) {
	key, err := noise.DH25519.GenerateKeypair(rand.Reader)
	if err != nil {
		panic("key generation failure")
	}
	switch keyType {
	case KeyTypeEphemeral:
		state.noiseConfig.StaticKeypair = key
	case KeyTypeStatic:
		state.noiseConfig.StaticKeypair = key

	}
	log.Debugf("[%v] generated new key %v of type '%v'", state.name, key, keyType)
}
