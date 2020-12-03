package model

import (
	"reflect"
	"testing"
)

func TestAddImage(t *testing.T) {
	var tests = []struct{
		tag Tag
		image Image
		expected []string
	}{
		{
			tag: Tag{
				Name: "empty",
				Images: []string{},
			},
			image: Image{
				ID: "my-image",
			},
			expected: []string{"my-image"},
		},
		{
			tag: Tag{
				Name: "not-empty",
				Images: []string{"my-image"},
			},
			image: Image{
				ID: "my-new-image",
			},
			expected: []string{"my-image", "my-new-image"},
		},
		{
			tag: Tag{
				Name: "duplicate",
				Images: []string{"my-image"},
			},
			image: Image{
				ID: "my-image",
			},
			expected: []string{"my-image"},
		},
	}

	for _, test := range tests {
		t.Run(
			test.tag.Name,
			func(t *testing.T) {
				test.tag.AddImage(&test.image)
				if !reflect.DeepEqual(test.tag.Images, test.expected) {
					t.Errorf("Expected %s, got %s", test.expected, test.tag.Images)
				}
			},
		)
	}
}
