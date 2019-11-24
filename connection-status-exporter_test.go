package main

import (
	"testing"
)

func TestSocketCheck(test *testing.T) {
	cases := []struct {
		caseName       string
		socketToTest   socket
		expectedToFail bool
	}{
		{
			caseName: "Empty socket name",
			socketToTest: socket{
				Name: "",
			},
			expectedToFail: true,
		},
		{
			caseName: "Empty host",
			socketToTest: socket{
				Name: "Test socket",
			},
			expectedToFail: true,
		},
		{
			caseName: "Empty port",
			socketToTest: socket{
				Name: "Test socket",
				Host: "localhost",
			},
			expectedToFail: true,
		},
		{
			caseName: "Incorrect protocol",
			socketToTest: socket{
				Name:     "Test socket",
				Host:     "localhost",
				Port:     80,
				Protocol: "incorrectProtocol",
			},
			expectedToFail: true,
		},
		{
			caseName: "Correct socket",
			socketToTest: socket{
				Name: "Test socket",
				Host: "localhost",
				Port: 80,
			},
			expectedToFail: false,
		},
	}

	for _, thisCase := range cases {
		output := thisCase.socketToTest.check()

		// For every case, checks if fails when it should fail
		// and goes ok when it should go ok
		if ((output != nil) && (thisCase.expectedToFail == false)) ||
			((output == nil) && (thisCase.expectedToFail == true)) {
			test.Errorf("Test not passed: %s", thisCase.caseName)
		}

	}
}
