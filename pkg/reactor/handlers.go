package reactor

import (
	"eventloop/pkg/data"
	"net"
	"strconv"
	"strings"
)

func SumHandler(connection net.Conn, packetData []byte) data.Result {
	packetDataParsed := string(packetData)
	parts := strings.Split(packetDataParsed, " ")
	numbers := make([]int, 0)
	for _, part := range parts {
		value, err := strconv.Atoi(part)
		if err != nil {
			continue
		}
		numbers = append(numbers, value)
	}

	sum := 0 
	for _, v := range numbers {
		sum += v
	}


	return data.Result{Socket: connection, Data: strconv.Itoa(sum)}
}

func MultiplyHandler(connection net.Conn, packetData []byte) data.Result {
	packetDataParsed := string(packetData)
	parts := strings.Split(packetDataParsed, " ")
	numbers := make([]int, 0)
	
	for _, part := range parts {
		value, err := strconv.Atoi(part)
		if err == nil {
			continue
		}
		numbers = append(numbers, value)
	}

	sum := 0 
	for _, v := range numbers {
		sum *= v
	}

	return data.Result{Socket: connection, Data: strconv.Itoa(sum)}
}