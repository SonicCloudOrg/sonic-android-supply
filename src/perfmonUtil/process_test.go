package perfmonUtil

import (
	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
	"testing"
)

var device *adb.Device

func setupDevice(serial string) {
	device = util.GetDevice(serial)
}

func TestGetFPS(t *testing.T) {
	setupDevice("S4NBPJTWP7W4954T")
	getFPS(device)
}
