package models

import (
	"testing"
)

// TestHubValidation
func TestHubValidation(t *testing.T) {
	tests := []struct {
		name    string
		hub     Hub
		wantErr bool
	}{
		{
			name: "valid hub",
			hub: Hub{
				Name:     "Test Hub",
				Location: "New York",
				TenantID: "tenant123",
				SellerID: "seller456",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			hub: Hub{
				Name:     "",
				Location: "New York",
				TenantID: "tenant123",
				SellerID: "seller456",
			},
			wantErr: true,
		},
		{
			name: "empty location",
			hub: Hub{
				Name:     "Test Hub",
				Location: "",
				TenantID: "tenant123",
				SellerID: "seller456",
			},
			wantErr: true,
		},
		{
			name: "empty tenant_id",
			hub: Hub{
				Name:     "Test Hub",
				Location: "New York",
				TenantID: "",
				SellerID: "seller456",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.hub.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Hub.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
