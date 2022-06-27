package port

import (
	"errors"
	"math/rand"
	"net"
	"strconv"
	"strings"
)

type Chooser interface {
	Choose() (uint16, error)
}

type randomPortChooser struct {
	baseAddress *net.TCPAddr
}

func (p *randomPortChooser) Choose() (uint16, error) {
	addr, err := net.ListenTCP("tcp", p.baseAddress)
	if err != nil {
		return 0, err
	}
	defer addr.Close()
	tcpAddress, ok := addr.Addr().(*net.TCPAddr)
	if !ok {
		return 0, errors.New("address was an invalid TCP Address")
	}
	return uint16(tcpAddress.Port), nil
}

func NewRandomPortChooser() (*randomPortChooser, error) {
	randomPortAddress, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		return nil, err
	}
	return &randomPortChooser{baseAddress: randomPortAddress}, nil
}

type rangePortChooser struct {
	portRange     []uint16
	selectedPorts map[uint16]struct{}
}

func (r *rangePortChooser) Choose() (uint16, error) {
	if len(r.selectedPorts) == len(r.portRange) {
		r.selectedPorts = map[uint16]struct{}{}
	}
	chosenPort := r.portRange[rand.Intn(len(r.portRange))]
	for r.isChosen(chosenPort) {
		chosenPort = r.portRange[rand.Intn(len(r.portRange))]
	}
	r.selectedPorts[chosenPort] = struct{}{}
	return chosenPort, nil
}

func (r *rangePortChooser) isChosen(port uint16) bool {
	_, ok := r.selectedPorts[port]
	return ok
}

func (r *rangePortChooser) DeselectPort(port uint16) {
	delete(r.selectedPorts, port)
}

func NewRangePortChooser(portRangeStr string) (*rangePortChooser, error) {
	portRange, err := getRangeFromStr(portRangeStr)
	if err != nil {
		return nil, err
	}
	return &rangePortChooser{portRange: portRange, selectedPorts: map[uint16]struct{}{}}, nil
}

func getRangeFromStr(rangeStr string) ([]uint16, error) {
	rangeNumStrs := strings.Split(rangeStr, "-")
	if len(rangeNumStrs) != 2 {
		return nil, errors.New("range can only have two numbers")
	}
	rangeNums := make([]uint16, 2)
	for i := 0; i < 2; i++ {
		portUint64, err := strconv.Atoi(strings.Trim(rangeNumStrs[i], " "))
		if err != nil {
			return nil, err
		}
		rangeNums[i] = uint16(portUint64)
	}

	if rangeNums[0] >= rangeNums[1] {
		return nil, errors.New("left number in range must be strictly lower than the right number")
	}
	portRange := []uint16{}
	for port := rangeNums[0]; port <= rangeNums[1]; port++ {
		portRange = append(portRange, port)
	}
	return portRange, nil
}
