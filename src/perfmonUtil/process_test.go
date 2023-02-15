package perfmonUtil

import (
	"fmt"
	"testing"

	"github.com/SonicCloudOrg/sonic-android-supply/src/adb"
	"github.com/SonicCloudOrg/sonic-android-supply/src/util"
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

func TestGetCurrentActivity(t *testing.T) {
	setupDevice("91cf5f1c")
	fmt.Println(GetCurrentActivity(device))
}

func TestGet(t *testing.T) {
	setupDevice("91cf5f1c")
	fmt.Println(GetPidOnPackageName(device, "com.tencent.mm"))
}
