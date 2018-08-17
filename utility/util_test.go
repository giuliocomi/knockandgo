package utility

import "testing"

func TestContainsPort(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   []int
		param2   int
		expected bool
	}{
		{[]int{1, 22, 3, 5050, 80, 8080, 445}, 445, true},
		{[]int{1, 22, 3, 5050, 80, 8080, 445}, -22, false},
		{[]int{1, 22, 3, 5050, 80, 8080, 445}, 12321, false},
		{[]int{1, 22, 3, 5050, 80, 8080, 445}, 808080, false},
		{[]int{1, 22, 22, 22, 80, 8080, 445}, 22, true},
		{[]int{1, 22, 3, 5050, 80, 8080, 445000}, 445000, false},
	}
	for _, test := range tests {
		actual := ContainsPort(test.param1, test.param2)
		if actual != test.expected {
			t.Errorf("Expected ContainsPort(%q, %q) to be %v, got %v", test.param1, test.param2, test.expected, actual)
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
			t.Errorf("Expected IsValidIP4(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}

func SliceAtoi(slice_array[] string)([] int, error) {
    slice_of_ints := make([] int, 0, len(slice_array))
    for _,
    a := range slice_array {
        integer, err := strconv.Atoi(a)
        if err != nil {
            return slice_of_ints, err
        }
        slice_of_ints = append(slice_of_ints, integer)
    }
    return slice_of_ints, nil
}

func TestSliceAtoi(t *testing.T) {
	t.Parallel()
	
	var tests = []struct {
		param    []string
		expected []int
		err	 error
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
		actual, err := IsValidIP4(test.param)
		if err != nil {
			t.Errorf("Expected SliceAtoi(%q) to return %q, got %v", test.param, test.expected, err)
		}
	}
}
