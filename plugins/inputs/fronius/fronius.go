package fronius

import (
	"net/http"
	"time"

	"github.com/KalleDK/go-fronius-solar/fronius"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

const (
	measurementDevice = "fronius_device"
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

	phaseData, err := client.Get3PInverterData()
	if err != nil {
		return err
	}

	commonData, err := client.GetCommonInverterData()
	if err != nil {
		return err
	}

	timestamp := time.Now()

	acc.AddFields(
		measurementDevice,
		map[string]interface{}{
			"p1_current":  float64(phaseData.Phase1.Current) / float64(fronius.Ampere),
			"p1_voltage":  float64(phaseData.Phase1.Voltage) / float64(fronius.Volt),
			"p2_current":  float64(phaseData.Phase2.Current) / float64(fronius.Ampere),
			"p2_voltage":  float64(phaseData.Phase2.Voltage) / float64(fronius.Volt),
			"p3_current":  float64(phaseData.Phase3.Current) / float64(fronius.Ampere),
			"p3_voltage":  float64(phaseData.Phase3.Voltage) / float64(fronius.Volt),
			"accurrent":   float64(commonData.ACCurrent) / float64(fronius.Ampere),
			"acvoltage":   float64(commonData.ACVoltage) / float64(fronius.Volt),
			"aceffect":    float64(commonData.ACEffect) / float64(fronius.Watt),
			"dccurrent":   float64(commonData.DCCurrent) / float64(fronius.Ampere),
			"dcvoltage":   float64(commonData.DCVoltage) / float64(fronius.Volt),
			"dayenergy":   float64(commonData.DayEnergy) / float64(fronius.Watthour),
			"yearenergy":  float64(commonData.YearEnergy) / float64(fronius.Watthour),
			"totalenergy": float64(commonData.TotalEnergy) / float64(fronius.Watthour),
			"frequency":   float64(commonData.Frequency) / float64(fronius.Hertz),
		},
		nil,
		timestamp)

	return nil
}

func (s *Fronius) Gather(acc telegraf.Accumulator) error {
	for _, device := range s.Devices {
		if err := s.gather(device, acc); err != nil {
			acc.AddError(err)
		}
	}
	return nil
}

func init() {
	inputs.Add("fronius", func() telegraf.Input { return &Fronius{} })
}
