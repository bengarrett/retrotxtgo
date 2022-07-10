package example_test

// func TestPrint(t *testing.T) {
// 	color.Enable = false
// 	tests := []struct {
// 		name string
// 		tmpl string
// 		want string
// 	}{
// 		{"empty", "", ""},
// 		{"word", "Hello", "Hello"},
// 		{"words", "Hello world", "Hello world"},
// 		{"comment", "Hello # world", "Hello # world"},
// 		{"comments", "Hello # hash # world", "Hello # hash # world"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := example.Print(tt.tmpl); strings.TrimSpace(got) != tt.want {
// 				t.Errorf("Print() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
