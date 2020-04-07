package fronius

import (
	"net/http"
	"time"

	"github.com/KalleDK/go-fronius-solar/fronius"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type Fronius struct {
	Devices []string `toml:"devices"`
}

func (s *Fronius) Description() string {
	return "The Fronius Solar API is a means for third parties to obtain data from various Fronius devices (inverters, SensorCards, StringControls)"
}

func (s *Fronius) SampleConfig() string {
	return `
  ## Indicate if everything is fine
  devices = ["http://192.168.255.13"]
`
}

func (s *Fronius) Init() error {
	return nil
}

func (s *Fronius) gather(device string, acc telegraf.Accumulator) error {
	client := fronius.Client{
		HttpClient: http.DefaultClient,
		BaseUrl:    device,
	}
	data, err := client.Get3PInverterData()
	if err != nil {
		return err
	}

	timestamp := time.Now()

	acc.AddFields(
		"phase4",
		map[string]interface{}{
			"current": float64(data.Phase1.Current) / float64(fronius.Ampere),
			"voltage": float64(data.Phase1.Voltage) / float64(fronius.Volt),
		},
		nil,
		timestamp)

	return nil
}

func (s *Fronius) Gather(acc telegraf.Accumulator) error {
	for _, device := range s.Devices {
		if err := s.gather(device, acc); err != nil {
			return err
		}
	}
	return nil
}

func init() {
	inputs.Add("fronius", func() telegraf.Input { return &Fronius{} })
}
