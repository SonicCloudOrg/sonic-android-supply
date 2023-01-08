package perfmonUtil

import (
	"fmt"
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
	r, _ := getProcessFPSBySurfaceFlinger(device, "com.android.browser")
	fmt.Println(r)
}
