package globalconfig

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	FeatureGroupChat         = "group_chat"
	FeatureVoiceMessage      = "voice_message"
	FeatureFileUpload        = "file_upload"
	FeatureVideoMessage      = "video_message"
	FeatureEnableRTC         = "enable_rtc"
	FeatureAllowOpenRegister = "allow_open_register"
	FeatureAllowCreateGroup  = "allow_create_group"
	FeatureAllowAddFriend    = "allow_add_friend"

	LimitMaxGroupMembers  = "max_group_members"
	LimitMaxMessageSizeMB = "max_message_size_mb"
	LimitMaxFileSizeMB    = "max_file_size_mb"
	LimitMaxDAU           = "max_dau"
	LimitMaxStorageGB     = "max_storage_gb"
	LimitLoginIPWhitelist = "login_ip_whitelist"
)

type Config struct {
	ServerID      string         `json:"server_id"`
	ConfigVersion string         `json:"config_version"`
	IssuedAt      int64          `json:"issued_at"`
	ExpiresAt     int64          `json:"expires_at"`
	Features      map[string]any `json:"features"`
	Limits        map[string]any `json:"limits"`
}

type ConfigPackage struct {
	ServerID      string         `json:"server_id"`
	ConfigVersion string         `json:"config_version"`
	IssuedAt      int64          `json:"issued_at"`
	ExpiresAt     int64          `json:"expires_at"`
	Features      map[string]any `json:"features"`
	Limits        map[string]any `json:"limits"`
	RawJSON       string         `json:"-"`
}

type Reader struct {
	features map[string]any
	limits   map[string]any
}

func NewReader(features map[string]any, limits map[string]any) Reader {
	return Reader{
		features: features,
		limits:   limits,
	}
}

func NewReaderFromJSON(configJSON string) (Reader, error) {
	if strings.TrimSpace(configJSON) == "" {
		return Reader{}, fmt.Errorf("global config is empty")
	}

	var config Config
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return Reader{}, fmt.Errorf("parse global config failed: %w", err)
	}
	return NewReader(config.Features, config.Limits), nil
}

func (r Reader) FeatureEnabled(key string) bool {
	value, exists, err := r.LookupBool(key)
	if !exists || err != nil {
		return defaultBool(key)
	}
	return value
}

func (r Reader) LimitInt64(key string) int64 {
	value, exists, err := r.LookupInt64(key)
	if !exists || err != nil {
		return defaultInt64(key)
	}
	return value
}

func (r Reader) LimitStringSlice(key string) []string {
	value, exists, err := r.LookupStringSlice(key)
	if !exists || err != nil {
		return defaultStringSlice(key)
	}
	return value
}

func (r Reader) LookupBool(key string) (value bool, exists bool, err error) {
	if len(r.features) == 0 {
		return false, false, nil
	}
	raw, ok := r.features[key]
	if !ok {
		return false, false, nil
	}
	value, err = toBool(raw)
	return value, true, err
}

func (r Reader) LookupInt64(key string) (value int64, exists bool, err error) {
	if len(r.limits) == 0 {
		return 0, false, nil
	}
	raw, ok := r.limits[key]
	if !ok {
		return 0, false, nil
	}
	value, err = toInt64(raw)
	return value, true, err
}

func (r Reader) LookupStringSlice(key string) (value []string, exists bool, err error) {
	if len(r.limits) == 0 {
		return nil, false, nil
	}
	raw, ok := r.limits[key]
	if !ok {
		return nil, false, nil
	}
	value, err = toStringSlice(raw)
	return value, true, err
}

func (r Reader) AllowOpenRegister() bool {
	return r.FeatureEnabled(FeatureAllowOpenRegister)
}

func (r Reader) AllowCreateGroup() bool {
	return r.FeatureEnabled(FeatureAllowCreateGroup)
}

func (r Reader) AllowAddFriend() bool {
	return r.FeatureEnabled(FeatureAllowAddFriend)
}

func (r Reader) FileUploadEnabled() bool {
	return r.FeatureEnabled(FeatureFileUpload)
}

func (r Reader) MaxGroupMembers() int64 {
	return r.LimitInt64(LimitMaxGroupMembers)
}

func (r Reader) MaxMessageSizeMB() int64 {
	return r.LimitInt64(LimitMaxMessageSizeMB)
}

func (r Reader) MaxFileSizeMB() int64 {
	return r.LimitInt64(LimitMaxFileSizeMB)
}

func (r Reader) MaxDAU() int64 {
	return r.LimitInt64(LimitMaxDAU)
}

func (r Reader) MaxStorageGB() int64 {
	return r.LimitInt64(LimitMaxStorageGB)
}

func (r Reader) LoginIPWhitelist() []string {
	return r.LimitStringSlice(LimitLoginIPWhitelist)
}

func defaultBool(string) bool {
	return false
}

func defaultInt64(string) int64 {
	return 0
}

func defaultStringSlice(string) []string {
	return []string{}
}

func toBool(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(v))
		if err != nil {
			return false, fmt.Errorf("invalid bool string %q", v)
		}
		return parsed, nil
	case float64:
		return numericBool(v)
	case float32:
		return numericBool(float64(v))
	case int:
		return intBool(int64(v))
	case int8:
		return intBool(int64(v))
	case int16:
		return intBool(int64(v))
	case int32:
		return intBool(int64(v))
	case int64:
		return intBool(v)
	case uint:
		return uintBool(uint64(v))
	case uint8:
		return uintBool(uint64(v))
	case uint16:
		return uintBool(uint64(v))
	case uint32:
		return uintBool(uint64(v))
	case uint64:
		return uintBool(v)
	case json.Number:
		i, err := v.Int64()
		if err != nil {
			return false, fmt.Errorf("invalid bool number %q", v)
		}
		return intBool(i)
	default:
		return false, fmt.Errorf("invalid bool type %T", value)
	}
}

func numericBool(value float64) (bool, error) {
	if value == 0 {
		return false, nil
	}
	if value == 1 {
		return true, nil
	}
	return false, fmt.Errorf("invalid bool number %v", value)
}

func intBool(value int64) (bool, error) {
	if value == 0 {
		return false, nil
	}
	if value == 1 {
		return true, nil
	}
	return false, fmt.Errorf("invalid bool number %d", value)
}

func uintBool(value uint64) (bool, error) {
	if value == 0 {
		return false, nil
	}
	if value == 1 {
		return true, nil
	}
	return false, fmt.Errorf("invalid bool number %d", value)
}

func toInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return uintToInt64(uint64(v))
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return uintToInt64(v)
	case float32:
		return floatToInt64(float64(v))
	case float64:
		return floatToInt64(v)
	case string:
		parsed, err := strconv.ParseInt(strings.TrimSpace(v), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid int64 string %q", v)
		}
		return parsed, nil
	case json.Number:
		parsed, err := v.Int64()
		if err != nil {
			return 0, fmt.Errorf("invalid int64 number %q", v)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("invalid int64 type %T", value)
	}
}

func uintToInt64(value uint64) (int64, error) {
	if value > math.MaxInt64 {
		return 0, fmt.Errorf("uint64 value %d overflows int64", value)
	}
	return int64(value), nil
}

func floatToInt64(value float64) (int64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("invalid int64 float %v", value)
	}
	if math.Trunc(value) != value {
		return 0, fmt.Errorf("non-integer float %v", value)
	}
	if value > math.MaxInt64 || value < math.MinInt64 {
		return 0, fmt.Errorf("float %v overflows int64", value)
	}
	return int64(value), nil
}

func toStringSlice(value any) ([]string, error) {
	switch v := value.(type) {
	case []string:
		return append([]string(nil), v...), nil
	case []any:
		result := make([]string, 0, len(v))
		for _, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("invalid string slice item type %T", item)
			}
			result = append(result, s)
		}
		return result, nil
	case string:
		if strings.TrimSpace(v) == "" {
			return []string{}, nil
		}
		parts := strings.Split(v, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				result = append(result, part)
			}
		}
		return result, nil
	default:
		return nil, fmt.Errorf("invalid string slice type %T", value)
	}
}
