package info_hub

import (
	"testing"
)

func TestInfoHub_Refresh(t *testing.T) {

	infoHub := NewInfoHub("")
	defer infoHub.Close()
	err := infoHub.Refresh()
	if err != nil {
		t.Fatal(err)
	}

	println(infoHub.SystemInfos.GetDriverVersion())
}

func BenchmarkInfoHub_Refresh(b *testing.B) {

	infoHub := NewInfoHub("")
	defer infoHub.Close()
	err := infoHub.Refresh()
	if err != nil {
		b.Fatal(err)
	}

	println(infoHub.SystemInfos.GetDriverVersion())
}
