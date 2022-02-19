package config_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"task_2_3_4/internal/pkg/config"
	"testing"
)

func Example() {
	cnfg, err := config.NewAppConfig()
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("%#v", cnfg)
	// Output:
	//&config.AppConfig{Path:".", Workers:5, DeleteDublicates:false, PrintResult:true}
}

func TestAppConfig_Check(t *testing.T) {
	var tests = []struct {
		input       config.AppConfig
		expectedErr error
	}{
		{
			input:       config.AppConfig{Path: ".", Workers: 5, DeleteDublicates: false, PrintResult: false},
			expectedErr: nil,
		},
		{
			input:       config.AppConfig{Path: "", Workers: 5, DeleteDublicates: true, PrintResult: false},
			expectedErr: fmt.Errorf("Path cant be empty"),
		},
		{
			input:       config.AppConfig{Path: "../files", Workers: 55, DeleteDublicates: false, PrintResult: false},
			expectedErr: fmt.Errorf("Amount of workers is limited from 1 to 50"),
		},
		{
			input:       config.AppConfig{Path: "../program", Workers: 0, DeleteDublicates: true, PrintResult: false},
			expectedErr: fmt.Errorf("Amount of workers is limited from 1 to 50"),
		},
	}

	for _, tt := range tests {
		cnfg, err := config.NewAppConfig()
		if err != nil {
			t.Errorf(err.Error())
		}

		cnfg = &tt.input
		assert.Equal(t, tt.expectedErr, cnfg.Validate(), "they should be equal")
	}
}
