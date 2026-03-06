package safety

import (
	"context"
	"testing"
)

func TestCheckURL_BlackListDomain(t *testing.T) {
	blackList := map[string]struct{}{
		"phishing.com":     {},
		"malware-site.org": {},
	}

	tests := []struct {
		name      string
		url       string
		wantSafe  bool
		wantLevel string
	}{
		{
			name:      "黑名单域名",
			url:       "https://phishing.com/fake-login",
			wantSafe:  false,
			wantLevel: "danger",
		},
		{
			name:      "正常域名",
			url:       "https://www.google.com",
			wantSafe:  true,
			wantLevel: "safe",
		},
		{
			name:      "另一个黑名单域名",
			url:       "https://malware-site.org/download",
			wantSafe:  false,
			wantLevel: "danger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckURL(context.Background(), tt.url, blackList, nil, "")
			if result.IsSafe != tt.wantSafe {
				t.Errorf("CheckURL() IsSafe = %v, want %v", result.IsSafe, tt.wantSafe)
			}
			if result.RiskLevel != tt.wantLevel {
				t.Errorf("CheckURL() RiskLevel = %v, want %v", result.RiskLevel, tt.wantLevel)
			}
		})
	}
}

func TestCheckURL_URLFeatures(t *testing.T) {
	emptyBlackList := map[string]struct{}{}

	tests := []struct {
		name      string
		url       string
		wantSafe  bool
		wantLevel string
	}{
		{
			name:      "IP直连",
			url:       "http://192.168.1.1/admin",
			wantSafe:  true,
			wantLevel: "warning",
		},
		{
			name:      "非标端口",
			url:       "http://example.com:8080/page",
			wantSafe:  true,
			wantLevel: "warning",
		},
		{
			name:      "标准HTTPS端口",
			url:       "https://example.com:443/page",
			wantSafe:  true,
			wantLevel: "safe",
		},
		{
			name:      "正常URL",
			url:       "https://www.example.com/page",
			wantSafe:  true,
			wantLevel: "safe",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckURL(context.Background(), tt.url, emptyBlackList, nil, "")
			if result.IsSafe != tt.wantSafe {
				t.Errorf("CheckURL() IsSafe = %v, want %v", result.IsSafe, tt.wantSafe)
			}
			if result.RiskLevel != tt.wantLevel {
				t.Errorf("CheckURL() RiskLevel = %v, want %v", result.RiskLevel, tt.wantLevel)
			}
		})
	}
}

func TestCheckURL_InvalidURL(t *testing.T) {
	emptyBlackList := map[string]struct{}{}
	result := CheckURL(context.Background(), "://invalid", emptyBlackList, nil, "")
	if result.IsSafe {
		t.Error("expected invalid URL to be unsafe")
	}
	if result.RiskLevel != "danger" {
		t.Errorf("expected danger, got %s", result.RiskLevel)
	}
}

func TestScoreToLevel(t *testing.T) {
	tests := []struct {
		score float64
		want  string
	}{
		{0.0, "safe"},
		{0.2, "safe"},
		{0.3, "warning"},
		{0.5, "warning"},
		{0.69, "warning"},
		{0.7, "danger"},
		{1.0, "danger"},
	}
	for _, tt := range tests {
		got := scoreToLevel(tt.score)
		if got != tt.want {
			t.Errorf("scoreToLevel(%f) = %s, want %s", tt.score, got, tt.want)
		}
	}
}
