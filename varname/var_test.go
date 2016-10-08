package varname
import "testing"
import "fmt"

func TestVar(t *testing.T) {
	for _, c := range []struct {
		itr int
		signs []rune 
		want []string
	}{
		{5, []rune("abcd"), []string{"a0", "b0", "c0", "d0", "a1"}},
		{10, []rune("abC"), []string{"a0", "b0", "C0", "a1", "b1", "C1", "a2", "b2", "C2", "a3"}},
		{7, []rune{},[]string{"0", "1", "2", "3", "4", "5", "6"}},
	} {
		t.Log(fmt.Sprintf("Running test with %v iterations and %v signs", c.itr, c.signs))
		v := NewVar(c.signs)
		got := RunIterationUsingGetAndInc(t, c.itr, &v)
		if !StringArrayEquals(got, c.want) {
			t.Errorf("Get(%q) == %q, want %q", c.signs, got, c.want)
		}
	}
}

func TestNextVar(t *testing.T) {
	for _, c := range []struct {
		itr int
		signs []rune 
		want []string
	}{
		{5, []rune("abcd"), []string{"a0", "b0", "c0", "d0", "a1"}},
		{10, []rune("abC"), []string{"a0", "b0", "C0", "a1", "b1", "C1", "a2", "b2", "C2", "a3"}},
		{7, []rune{},[]string{"0", "1", "2", "3", "4", "5", "6"}},
	} {
		t.Log(fmt.Sprintf("Running test with %v iterations and %v signs", c.itr, c.signs))
		v := NewVar(c.signs)
		got := RunIterationUsingNext(t, c.itr, &v)
		if !StringArrayEquals(got, c.want) {
			t.Errorf("Next(%q) == %q, want %q", c.signs, got, c.want)
		}
	}
}	

func StringArrayEquals(a []string, b []string) bool {
    if len(a) != len(b) {
        return false
    }
    for i, v := range a {
        if v != b[i] {
            return false
        }
    }
    return true
}

func RunIterationUsingNext(t *testing.T, itr int, v *Var) []string{
	var res []string
	for i := 0; i < itr; i++ {
		res = append(res, v.Next())
	}
	return res
}

func RunIterationUsingGetAndInc(t *testing.T, itr int, v *Var) []string{
	var res []string
	for i := 0; i < itr; i++ {
		r := v.Get()
		res = append(res, r)
		v.Inc()                 	
	}
	return res
}
