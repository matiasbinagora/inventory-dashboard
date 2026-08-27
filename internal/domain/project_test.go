package domain

import "testing"

func TestValidation(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
		link  PublicLink
	}{
		{name: "github public link", valid: true, link: PublicLink{Kind: GitHub, URL: "https://github.com/example/project"}},
		{name: "http rejected", link: PublicLink{Kind: GitHub, URL: "http://github.com/example/project"}},
		{name: "unsupported host rejected", link: PublicLink{Kind: GitHub, URL: "https://example.com/project"}},
		{name: "userinfo rejected", link: PublicLink{Kind: GitHub, URL: "https://user:pass@github.com/example/project"}},
		{name: "private query rejected", link: PublicLink{Kind: GitHub, URL: "https://github.com/example/project?token=secret"}},
		{name: "localhost rejected", link: PublicLink{Kind: GitHub, URL: "https://127.0.0.1/example/project"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.link.Validate()
			if (err == nil) != tt.valid {
				t.Fatalf("Validate() error=%v, valid=%v", err, tt.valid)
			}
		})
	}
}

func TestMediaValidationRequiresManagedRelativePath(t *testing.T) {
	valid := MediaAsset{Role: Screenshot, Source: "media/project/screenshot.png", Curated: true}
	if err := valid.Validate(); err != nil {
		t.Fatal(err)
	}
	for _, source := range []string{"/tmp/screenshot.png", "../secret.txt", "source/main.go", "media/project/.env", "media/project/transcript.txt", "media/project/build.zip"} {
		if err := (MediaAsset{Role: Screenshot, Source: source, Curated: true}).Validate(); err == nil {
			t.Fatalf("Validate(%q) succeeded", source)
		}
	}
	if err := (MediaAsset{Role: Screenshot, Source: "media/project/screenshot.png"}).Validate(); err == nil {
		t.Fatal("uncurated media was accepted")
	}
	if err := (MediaAsset{Role: Video, Source: "https://localhost/demo.mp4", Curated: true}).Validate(); err == nil {
		t.Fatal("private video host was accepted")
	}
}

func TestProjectValidationRejectsSensitiveCuratedText(t *testing.T) {
	for _, description := range []string{"transcript of the meeting", "password=secret", "package main"} {
		if err := (Project{Name: "Safe name", Description: description}).Validate(); err == nil {
			t.Fatalf("sensitive description %q was accepted", description)
		}
	}
}
