package model

import (
	"strconv"
	"testing"

	"github.com/spf13/viper"
)

func TestGetFilename(t *testing.T) {
	var tests = []struct{
		image Image
		expected string
	}{
		{
			image: Image{
				ID: "abcd",
				Extension: "gif",
			},
			expected: "abcd.gif",
		},
		{
			image: Image{
				ID: "test",
				Extension: "webm",
			},
			expected: "test.webm",
		},
	}

	for index, test := range tests {
		t.Run(
			strconv.Itoa(index),
			func(t *testing.T) {
				actual := test.image.GetFilename()
				if actual != test.expected {
					t.Errorf("Got %s, expected %s", actual, test.expected)
				}
			},
		)
	}
}

func TestGetPublicURL(t *testing.T) {
	var tests = []struct{
		image Image
		host string
		expected string
	}{
		{
			image: Image{
				ID: "abcd",
				Extension: "gif",
			},
			host: "http://gifs.example.com",
			expected: "http://gifs.example.com/files/abcd.gif",
		},
		{
			image: Image{
				ID: "test",
				Extension: "webm",
			},
			host: "https://gifs.senketsu.test/",
			expected: "https://gifs.senketsu.test/files/test.webm",
		},
		{
			image: Image{
				ID: "test2",
				Extension: "webm",
			},
			host: "https://gifs.senketsu.test/mysubfolder",
			expected: "https://gifs.senketsu.test/mysubfolder/files/test2.webm",
		},
	}

	for index, test := range tests {
		t.Run(
			strconv.Itoa(index),
			func(t *testing.T) {
				viper.Set("host", test.host)
				actual := test.image.GetPublicURL()
				if actual != test.expected {
					t.Errorf("Got %s, expected %s", actual, test.expected)
				}
			},
		)
	}
}
