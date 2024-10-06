package output

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/dayvillefire/pocsag-monitor/obj"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func init() {
	RegisterOutput("mqtt", func() Output { return &MQTTOutput{} })
}

type MQTTOutput struct {
	clientID string
	uri      string

	client mqtt.Client
}

func (m *MQTTOutput) Init(v string) error {
	var err error
	m.uri = v
	m.clientID, err = m.getLocalIP()
	if err != nil {
		return err
	}
	//log.Printf("clientid = %s", m.clientID)
	return m.connect()
}

func (m *MQTTOutput) SendMessage(a obj.AlphaMessage, channel, msg string) (string, error) {
	log.Printf("MQTTOutput.SendMessage(%#v, %s, %s)", a, channel, msg)
	//log.Printf("%#v", m.client)
	v := m.client.Publish(channel, 0, false, msg)
	//go func(mt mqtt.Token) {
	<-v.Done()
	if v.Error() != nil {
		log.Printf("ERR: MQTT.SendMessage: %s", v.Error().Error())
		return "", v.Error()
	}
	//}(v)
	return "", nil
}

func (m *MQTTOutput) connect() error {
	log.Printf("MQTTOutput.connect: uri = %s, clientId = %s", m.uri, m.clientID)
	uri, err := url.Parse(m.uri)
	log.Printf("MQTTOutput.connect: parsed uri = %#v", uri)
	if err != nil {
		log.Printf("ERR: %s", err.Error())
		return err
	}
	opts := m.createClientOptions(m.clientID, uri)
	m.client = mqtt.NewClient(opts)
	token := m.client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Printf("ERR: %s", err.Error())
		return err
	}
	return nil
}

func (m *MQTTOutput) createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", uri.Host))
	opts.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	opts.SetPassword(password)
	opts.SetClientID(clientId)
	return opts
}

func (m MQTTOutput) getLocalIP() (string, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			continue
		}
		for _, a := range aa {
			ip := net.ParseIP(strings.Split(a.String(), "/")[0])
			if ip.IsLoopback() {
				continue
			}
			if !ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast() && !ip.IsLinkLocalMulticast() && !ip.IsLoopback() {

				//log.Printf("ip = %s", strings.Split(a.String(), "/")[0])
				return strings.Split(a.String(), "/")[0], nil
			}
		}
	}
	return "", fmt.Errorf("no IP address found")
}
