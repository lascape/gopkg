package authx

import (
	"reflect"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		kv      map[string]interface{}
		expired time.Duration
		wantErr bool
	}{
		{
			name:   "jwt success",
			method: AuthNameJwt,
			kv: map[string]interface{}{
				"username": "admin",
				"password": "admin",
			},
			expired: time.Minute * 30,
			wantErr: false,
		},
		{
			name:   "jwt expired token",
			method: AuthNameJwt,
			kv: map[string]interface{}{
				"username": "admin",
				"password": "admin",
			},
			expired: -time.Minute * 30,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetAuth(tt.method).GenerateToken(tt.kv, tt.expired)
			if err != nil {
				t.Errorf("GenerateToken() error = %v", err)
				return
			}
			newKv, err := GetAuth(tt.method).ValidateToken(token)
			if tt.wantErr && err != nil {
				return
			}
			if err != nil {
				t.Errorf("GenerateToken() error = %v", err)
				return
			}
			if !reflect.DeepEqual(newKv, tt.kv) {
				t.Errorf("GenerateToken() got = %v, want %v", newKv, tt.kv)
				return
			}
		})
	}
}
