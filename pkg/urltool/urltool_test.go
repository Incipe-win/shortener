package urltool

import "testing"

func TestGetBasePath(t *testing.T) {
	type args struct {
		targetUrl string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "基本示例", args: args{targetUrl: "https://www.liwenzhou.com/posts/Go/golang-menu/"}, want: "golang-menu", wantErr: false},
		{name: "相对路径url", args: args{targetUrl: "/xxxx/1123"}, want: "", wantErr: true},
		{name: "空字符串", args: args{targetUrl: ""}, want: "", wantErr: true},
		{name: "带query的url", args: args{targetUrl: "https://www.liwenzhou.com/posts/Go/golang-menu/?a=1&b=2"}, want: "golang-menu", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBasePath(tt.args.targetUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBasePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCircularURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		shortDomain string
		want        bool
	}{
		{
			name:        "精确匹配短域名",
			inputURL:    "https://shortener.com/abc123",
			shortDomain: "shortener.com",
			want:        true,
		},
		{
			name:        "子域名匹配",
			inputURL:    "https://www.shortener.com/abc123",
			shortDomain: "shortener.com",
			want:        true,
		},
		{
			name:        "不同域名",
			inputURL:    "https://www.example.com/page",
			shortDomain: "shortener.com",
			want:        false,
		},
		{
			name:        "部分匹配但不是子域名",
			inputURL:    "https://notshortener.com/page",
			shortDomain: "shortener.com",
			want:        false,
		},
		{
			name:        "大小写不敏感",
			inputURL:    "https://SHORTENER.COM/abc",
			shortDomain: "shortener.com",
			want:        true,
		},
		{
			name:        "无效URL",
			inputURL:    "://invalid",
			shortDomain: "shortener.com",
			want:        false,
		},
		{
			name:        "空URL",
			inputURL:    "",
			shortDomain: "shortener.com",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsCircularURL(tt.inputURL, tt.shortDomain)
			if got != tt.want {
				t.Errorf("IsCircularURL(%q, %q) = %v, want %v", tt.inputURL, tt.shortDomain, got, tt.want)
			}
		})
	}
}
