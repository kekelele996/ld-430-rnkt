package config

import "testing"

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		setEnv  map[string]string
		wantErr bool
	}{
		{
			name: "development valid",
			setEnv: map[string]string{
				"APP_ENV":     "development",
				"JWT_SECRET":  "0123456789abcdef0123456789abcdef",
				"MONGO_URI":   "mongodb://localhost:27017",
				"REDIS_URL":   "redis://localhost:6379/0",
				"SERVER_PORT": "8080",
			},
		},
		{
			name: "development short jwt rejected",
			setEnv: map[string]string{
				"APP_ENV":    "development",
				"JWT_SECRET": "short",
				"MONGO_URI":  "mongodb://localhost:27017",
				"REDIS_URL":  "redis://localhost:6379/0",
			},
			wantErr: true,
		},
		{
			name: "production default jwt rejected",
			setEnv: map[string]string{
				"APP_ENV":              "production",
				"JWT_SECRET":           "change_me_to_a_long_random_secret",
				"MINIO_ACCESS_KEY":     "assethub-access",
				"MINIO_SECRET_KEY":     "assethub-secret-change-me-8d2e",
				"CORS_ALLOWED_ORIGINS": "http://localhost:19310",
				"ADMIN_PASSWORD":       "AssethubAdmin2026",
			},
			wantErr: true,
		},
		{
			name: "production wildcard cors rejected",
			setEnv: map[string]string{
				"APP_ENV":              "production",
				"JWT_SECRET":           "0123456789abcdef0123456789abcdef",
				"MINIO_ACCESS_KEY":     "assethub-access",
				"MINIO_SECRET_KEY":     "assethub-secret-change-me-8d2e",
				"CORS_ALLOWED_ORIGINS": "*",
				"ADMIN_PASSWORD":       "AssethubAdmin2026",
			},
			wantErr: true,
		},
		{
			name: "production valid",
			setEnv: map[string]string{
				"APP_ENV":              "production",
				"JWT_SECRET":           "0123456789abcdef0123456789abcdef",
				"MINIO_ACCESS_KEY":     "assethub-access",
				"MINIO_SECRET_KEY":     "assethub-secret-change-me-8d2e",
				"CORS_ALLOWED_ORIGINS": "http://localhost:19310",
				"ADMIN_PASSWORD":       "AssethubAdmin2026",
				"MONGO_URI":            "mongodb://mongo:27017",
				"REDIS_URL":            "redis://redis:6379/0",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.setEnv {
				t.Setenv(k, v)
			}
			_, err := Load()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
