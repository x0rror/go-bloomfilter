package config

import (
	"testing"
	"time"
)

func TestFactoryConfig_Validate(t *testing.T) {
	type fields struct {
		FilterConfig  FilterConfig
		RedisConfig   RedisConfig
		RotatorConfig RotatorConfig
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid: filter",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						BitmapTypeInMemory,
					},
					M: 100,
					K: 2,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid: filter M",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						BitmapTypeInMemory,
					},
					M: 0,
					K: 2,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: bitmap type",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						Type: "",
					},
					M: 100,
					K: 2,
				},
			},
			wantErr: true,
		},
		{
			name: "valid: redis",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						BitmapTypeRedis,
					},
					M: 100,
					K: 2,
				},
				RedisConfig: RedisConfig{
					Addr:    "localhost:6379",
					Timeout: 5 * time.Second,
					Key:     "filter-redis",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid: redis, empty addr",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						BitmapTypeRedis,
					},
					M: 100,
					K: 2,
				},
				RedisConfig: RedisConfig{
					Addr:    "",
					Timeout: 5 * time.Second,
					Key:     "filter-redis",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: redis, zero timeout",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						BitmapTypeRedis,
					},
					M: 100,
					K: 2,
				},
				RedisConfig: RedisConfig{
					Addr:    "localhost",
					Timeout: 0,
					Key:     "filter-redis",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: redis, empty key",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						BitmapTypeRedis,
					},
					M: 100,
					K: 2,
				},
				RedisConfig: RedisConfig{
					Addr:    "localhost",
					Timeout: 5 * time.Second,
					Key:     "",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid: rotator",
			fields: fields{
				FilterConfig: FilterConfig{
					BitmapConfig: BitmapConfig{
						BitmapTypeInMemory,
					},
					M: 100,
					K: 2,
				},
				RotatorConfig: RotatorConfig{
					Enable: true,
					Freq:   0,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fc := &FactoryConfig{
				FilterConfig:  tt.fields.FilterConfig,
				RedisConfig:   tt.fields.RedisConfig,
				RotatorConfig: tt.fields.RotatorConfig,
			}
			if err := fc.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRotatorConfig_Validate(t *testing.T) {
	type fields struct {
		Enable bool
		Freq   time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid",
			fields: fields{
				Enable: false,
				Freq:   1,
			},
			wantErr: false,
		},
		{
			name: "invalid: zero freq",
			fields: fields{
				Enable: false,
				Freq:   0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := RotatorConfig{
				Enable: tt.fields.Enable,
				Freq:   tt.fields.Freq,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
