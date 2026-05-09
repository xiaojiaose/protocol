package globalconfig

import "testing"

func TestReaderFeatureEnabled(t *testing.T) {
	reader := NewReader(map[string]any{
		FeatureGroupChat:         true,
		FeatureFileUpload:        false,
		FeatureEnableRTC:         "true",
		FeatureAllowCreateGroup:  1,
		FeatureAllowOpenRegister: "0",
		FeatureAllowAddFriend:    2,
	}, nil)

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{name: "bool true", key: FeatureGroupChat, want: true},
		{name: "bool false", key: FeatureFileUpload, want: false},
		{name: "string true", key: FeatureEnableRTC, want: true},
		{name: "number true", key: FeatureAllowCreateGroup, want: true},
		{name: "string false", key: FeatureAllowOpenRegister, want: false},
		{name: "invalid returns default", key: FeatureAllowAddFriend, want: false},
		{name: "missing returns default", key: "missing", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reader.FeatureEnabled(tt.key); got != tt.want {
				t.Fatalf("FeatureEnabled() = %v, want %v", got, tt.want)
			}
		})
	}

	if _, exists, err := reader.LookupBool(FeatureAllowAddFriend); !exists || err == nil {
		t.Fatalf("LookupBool(invalid) exists=%v err=%v, want exists and error", exists, err)
	}
}

func TestReaderLimitInt64(t *testing.T) {
	reader := NewReader(nil, map[string]any{
		LimitMaxGroupMembers:  float64(300),
		LimitMaxMessageSizeMB: "5",
		LimitMaxFileSizeMB:    int64(100),
		LimitMaxDAU:           float64(12.5),
	})

	tests := []struct {
		name string
		key  string
		want int64
	}{
		{name: "json number", key: LimitMaxGroupMembers, want: 300},
		{name: "string number", key: LimitMaxMessageSizeMB, want: 5},
		{name: "int64", key: LimitMaxFileSizeMB, want: 100},
		{name: "invalid returns default", key: LimitMaxDAU, want: 0},
		{name: "missing returns default", key: "missing", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reader.LimitInt64(tt.key); got != tt.want {
				t.Fatalf("LimitInt64() = %d, want %d", got, tt.want)
			}
		})
	}

	if _, exists, err := reader.LookupInt64(LimitMaxDAU); !exists || err == nil {
		t.Fatalf("LookupInt64(invalid) exists=%v err=%v, want exists and error", exists, err)
	}
}

func TestReaderLimitStringSlice(t *testing.T) {
	reader := NewReader(nil, map[string]any{
		"array_any":             []any{"192.168.1.1", "10.0.0.1"},
		"array_string":          []string{"127.0.0.1"},
		"comma_string":          "192.168.1.1, 10.0.0.1,,",
		LimitLoginIPWhitelist:   []any{"192.168.1.1", 10},
		LimitMaxMessageSizeMB:   100,
		LimitMaxGroupMembers:    []any{},
		LimitMaxFileSizeMB:      "",
		LimitMaxDAU:             " ",
		LimitMaxStorageGB:       "a,b",
		FeatureAllowAddFriend:   []string{"ignored"},
		FeatureAllowCreateGroup: []any{"ignored"},
	})

	tests := []struct {
		name string
		key  string
		want []string
	}{
		{name: "array any", key: "array_any", want: []string{"192.168.1.1", "10.0.0.1"}},
		{name: "array string", key: "array_string", want: []string{"127.0.0.1"}},
		{name: "comma string", key: "comma_string", want: []string{"192.168.1.1", "10.0.0.1"}},
		{name: "invalid returns default", key: LimitLoginIPWhitelist, want: []string{}},
		{name: "missing returns default", key: "missing", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reader.LimitStringSlice(tt.key)
			if len(got) != len(tt.want) {
				t.Fatalf("LimitStringSlice() len = %d, want %d (%v)", len(got), len(tt.want), got)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("LimitStringSlice()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}

	if _, exists, err := reader.LookupStringSlice(LimitLoginIPWhitelist); !exists || err == nil {
		t.Fatalf("LookupStringSlice(invalid) exists=%v err=%v, want exists and error", exists, err)
	}
}

func TestReaderSemanticMethodsAndFromJSON(t *testing.T) {
	reader, err := NewReaderFromJSON(`{"features":{"allow_open_register":true,"file_upload":true},"limits":{"max_file_size_mb":100,"login_ip_whitelist":"192.168.1.1,10.0.0.1"}}`)
	if err != nil {
		t.Fatal(err)
	}
	if !reader.AllowOpenRegister() {
		t.Fatal("AllowOpenRegister() = false, want true")
	}
	if !reader.FileUploadEnabled() {
		t.Fatal("FileUploadEnabled() = false, want true")
	}
	if got := reader.MaxFileSizeMB(); got != 100 {
		t.Fatalf("MaxFileSizeMB() = %d, want 100", got)
	}
	if got := reader.LoginIPWhitelist(); len(got) != 2 || got[0] != "192.168.1.1" || got[1] != "10.0.0.1" {
		t.Fatalf("LoginIPWhitelist() = %v, want two ips", got)
	}
}

func TestNewReaderFromJSONInvalid(t *testing.T) {
	if _, err := NewReaderFromJSON(""); err == nil {
		t.Fatal("NewReaderFromJSON(empty) succeeded, want error")
	}
	if _, err := NewReaderFromJSON("{"); err == nil {
		t.Fatal("NewReaderFromJSON(invalid) succeeded, want error")
	}
}
