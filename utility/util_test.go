package utility

import "testing"
import "bytes"
import "encoding/binary"

func TestContainsPort(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   []int
		param2   int
		expected bool
	}{
		{[]int{1, 445, 3, 5050, 80, 8080, 445}, 445, true},
		{[]int{1, 22, 3, 5050, 80, 8080, 445}, -22, false},
		{[]int{1, 22, 3, 5050, 80, 8080, 445}, 12321, false},
		{[]int{1, 22, 3, 5050, 80, 8080, 445}, 808080, false},
		{[]int{1, 22, 22, 22, 80, 8080, 445}, 22, true},
		{[]int{1, 22, 3, 5050, 80, 8080, 445000}, 445000, false},
	}
	for _, test := range tests {
		actual := ContainsPort(test.param1, test.param2)
		if actual != test.expected {
			t.Errorf("Expected ContainsPort(%v, %v) to be %v, got %v", test.param1, test.param2, test.expected, actual)
		}
	}
}

func TestIsValidIP4(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"127.0.0.1", true},
		{"", false},
		{"::", false},
		{"151.12.29.30", true},
		{"-151.12.29.30", false},
		{"127ssee.0.0.1", false},
		{"127.0.01", false},
		{"::1", false},
		{"localhost", true},
		{"287.212.12.2", false},
	}
	for _, test := range tests {
		actual := IsValidIP4(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsValidIP4(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func TestSliceAtoi(t *testing.T) {
	t.Parallel()

var tests = []struct {
		param    []string
		expected []int
	}{
		{[]string{"1","    22", "3", "asasaaaa", "80", "8080", "2445"}, []int{1, 22, 3, 80, 8080, 2445}},
		{[]string{"-1","  -  22", "3", "5050", "80", "8080", "2445"}, []int{-1, 3, 5050, 80, 8080, 2445}},
	}

	for _, test := range tests {
		actual, _ := SliceAtoi(test.param)
		
		var buf_expected bytes.Buffer
		binary.Write(&buf_expected, binary.BigEndian, test.expected)
		var buf_actual bytes.Buffer
		binary.Write(&buf_actual, binary.BigEndian, actual)

		if !bytes.Equal(buf_expected.Bytes(), buf_actual.Bytes()) {
			t.Errorf("Expected SliceAtoi(%v) to be %v, got %v", test.param, test.expected, actual)
		}
	}

}

