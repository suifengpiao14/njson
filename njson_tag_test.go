package njson

import (
	"fmt"
	"testing"
)

func TestNewNjsonTag(t *testing.T) {
	str := "fullname:_data._data.basic.customData.appointmentTime;format : json; formatPath:_data._data.basic.customData"
	njsonTag := NewNjsonTag(str)
	fmt.Println(njsonTag)
}
