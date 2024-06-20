package tflexregex

import "testing"

func validate(res bool, err error, onInvalid func()) {
	if err != nil || !res {
		onInvalid()
	}
}

func TestMatchAabc(t *testing.T) {
	res, err := Matches("aabc", []byte("aabc"))
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("aabc", []byte("abc"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })
}

func TestMatchStar(t *testing.T) {
	res, err := Matches("a*", []byte("aaaaaaaaaaa"))
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("a*", []byte{})
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("a*", []byte("a"))
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("ab*c*", []byte("abababccccccc"))
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("ab*", []byte("ababab"))
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("a*", []byte("b"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })

	res, err = Matches("a*", []byte("aaaab"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })

	res, err = Matches("ab*c*", []byte("abababcab"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })
}

func TestMatchPlus(t *testing.T) {
	res, err := Matches("a+", []byte("aaaaaaaaaaa"))
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("a+", []byte{})
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })

	res, err = Matches("a+", []byte("a"))
	validate(res, err, func() { t.Errorf("result is not positive %v: %v", res, err) })

	res, err = Matches("a+", []byte("b"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })

	res, err = Matches("a+", []byte("aaaab"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })

	res, err = Matches("ab+", []byte("aaaab"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })

	res, err = Matches("ab+", []byte("abbbbb"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })

	res, err = Matches("ab+", []byte("abababa"))
	validate(!res, err, func() { t.Errorf("result is not negative %v: %v", res, err) })
}
