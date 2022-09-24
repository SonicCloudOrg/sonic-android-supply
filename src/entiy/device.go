package entiy

type Device struct {
	Serial      string `json:"serial"`
	Status      string `json:"status"`
	Description string `json:"-"`
}

func (d *Device) String() string {
	return d.Serial
}
