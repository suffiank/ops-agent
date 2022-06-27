package port_test

import (
	"fmt"
	"testing"

	"github.com/GoogleCloudPlatform/ops-agent/internal/port"
)

func TestRandomPortChooser(t *testing.T) {
	t.Parallel()

	chooser, err := port.NewRandomPortChooser()
	if err != nil {
		t.Fatalf("couldn't make port chooser, got error: %v", err)
	}
	port, err := chooser.Choose()
	if err != nil {
		t.Fatalf("couldn't choose port, got error: %v", err)
	}
	t.Logf("Got port number: %d", port)
	if port == 0 {
		t.Fatalf("expected valid port number, got 0")
	}
}

func TestNewRangePortChooser(t *testing.T) {
	testCases := []struct {
		name     string
		rangeStr string
		success  bool
	}{
		{
			name:     "normal valid range",
			rangeStr: "8000 - 9000",
			success:  true,
		},
		{
			name:     "normal valid range, no spaces around dash",
			rangeStr: "5000-6000",
			success:  true,
		},
		{
			name:     "range no dash",
			rangeStr: "16 to 30",
		},
		{
			name:     "range too many dashes",
			rangeStr: "-40 - 50",
		},
		{
			name:     "range same numbers",
			rangeStr: "420 - 420",
		},
		{
			name:     "range left greater than right",
			rangeStr: "6000 - 5000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			chooser, err := port.NewRangePortChooser(tc.rangeStr)

			if tc.success {
				if err != nil {
					t.Fatalf("expected range string %q to construct successfully, got error: %v", tc.rangeStr, err)
				}
			} else {
				if err == nil {
					t.Log(chooser)
					t.Fatalf("expected range string %q not to construct successfully", tc.rangeStr)
				}
			}
		})
	}
}

func TestRangePortChooserChoosesPort(t *testing.T) {
	testCases := []struct {
		name  string
		start uint16
		end   uint16
	}{
		{
			name:  "regular short range",
			start: 20,
			end:   25,
		},
		{
			name:  "two number range",
			start: 300,
			end:   301,
		},
		{
			name:  "super long range",
			start: 3000,
			end:   4000,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rangeSize := int(tc.end-tc.start) + 1
			chooser, err := port.NewRangePortChooser(fmt.Sprintf("%d-%d", tc.start, tc.end))
			if err != nil {
				t.Fatalf("expected to successfully construct range chooser, got error: %v", err)
			}

			chosenPorts := map[uint16]struct{}{}
			for i := 0; i < rangeSize; i++ {
				port, _ := chooser.Choose()
				chosenPorts[port] = struct{}{}
			}

			if len(chosenPorts) != rangeSize {
				t.Log(chosenPorts)
				t.Fatalf("expected every unique port to be chosen in %d attempts", rangeSize)
			}
		})
	}
}
