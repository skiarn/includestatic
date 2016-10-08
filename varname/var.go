package varname 

import "fmt"
import "strconv"

func f2(a int) (int, int)	{ return a, a }
func f1(a, b int)		{ fmt.Println(a, a) }

func main() {
	var res []string
	v := NewVar([]rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"))
	for i := 0; i < 1024; i++ {
		fmt.Printf("%v - %s\n", i, v.Get())
		res = append(res, v.Get())
		v.Inc()	
	}
	fmt.Println(res)
}	

type Var struct {
	signs []rune 
	letterIndex int
	number int
	value string
}
func NewVar(signs []rune) Var {
	v := Var{}
	v.letterIndex = 0
	v.number = 0
	v.signs = signs
	v.update() 
	return v
}
//var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func(v *Var) Get() string {
	return v.value
}

//Next will increment and return next value.
func(v *Var) Next() string {
	v.Inc()
	return v.Get()
}
func(v *Var) String() string {
	return fmt.Sprintf("Letter I:%v N:%v Signs:%v\n", v.letterIndex, v.number, v.signs)
}
func(v *Var) Inc() {
	if len(v.signs) == 0 || v.letterIndex == len(v.signs)-1 {
		//if last
		v.number++ 
		v.letterIndex = 0
	} else {
		v.letterIndex++
	}
	
	//letter := string(v.signs[v.letterIndex:len(v.signs)-(len(v.signs)-v.letterLength)])
	v.update()
}

func(v *Var) update() {
	var letter string
	if v.letterIndex>=0 && v.letterIndex<len(v.signs) {
		letter = string(v.signs[v.letterIndex])
	}
	number := strconv.Itoa(v.number)
	v.value = letter+number
}
