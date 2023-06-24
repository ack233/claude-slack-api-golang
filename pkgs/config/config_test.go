package config

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init()

			fmt.Println(SlackConfig)
			fmt.Println(ViperConfig.GetString("logfile"))
			fmt.Println(ViperConfig.GetString("loglevel"))
			//if err := Init(); (err != nil) != tt.wantErr {
			//	t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			//}
		})
	}
}
