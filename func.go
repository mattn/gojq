package gojq

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
)

type function func(*env, *Func) func(interface{}) interface{}

var internalFuncs map[string]function

func init() {
	internalFuncs = map[string]function{
		"null":           noArgFunc(funcNull),
		"true":           noArgFunc(funcTrue),
		"false":          noArgFunc(funcFalse),
		"empty":          noArgFunc(funcEmpty),
		"length":         noArgFunc(funcLength),
		"utf8bytelength": noArgFunc(funcUtf8ByteLength),
		"keys":           noArgFunc(funcKeys),
		"has":            funcHas,
		"tonumber":       noArgFunc(funcToNumber),
		"type":           noArgFunc(funcType),
		"explode":        noArgFunc(funcExplode),
		"implode":        noArgFunc(funcImplode),
		"join":           funcJoin,
		"tojson":         noArgFunc(funcToJSON),
		"fromjson":       noArgFunc(funcFromJSON),
		"sin":            mathFunc("sin", math.Sin),
		"cos":            mathFunc("cos", math.Cos),
		"tan":            mathFunc("tan", math.Tan),
		"asin":           mathFunc("asin", math.Asin),
		"acos":           mathFunc("acos", math.Acos),
		"atan":           mathFunc("atan", math.Atan),
		"sinh":           mathFunc("sinh", math.Sinh),
		"cosh":           mathFunc("cosh", math.Cosh),
		"tanh":           mathFunc("tanh", math.Tanh),
		"asinh":          mathFunc("asinh", math.Asinh),
		"acosh":          mathFunc("acosh", math.Acosh),
		"atanh":          mathFunc("atanh", math.Atanh),
		"floor":          mathFunc("floor", math.Floor),
		"round":          mathFunc("round", math.Round),
		"ceil":           mathFunc("ceil", math.Ceil),
		"trunc":          mathFunc("trunc", math.Trunc),
		"fabs":           mathFunc("fabs", math.Abs),
		"sqrt":           mathFunc("sqrt", math.Sqrt),
		"cbrt":           mathFunc("cbrt", math.Cbrt),
		"exp":            mathFunc("exp", math.Exp),
		"exp10":          mathFunc("exp10", func(v float64) float64 { return math.Pow(10, v) }),
		"exp2":           mathFunc("exp2", math.Exp2),
		"expm1":          mathFunc("expm1", math.Expm1),
		"log":            mathFunc("log", math.Log),
		"log10":          mathFunc("log10", math.Log10),
		"log1p":          mathFunc("log1p", math.Log1p),
		"log2":           mathFunc("log2", math.Log2),
		"logb":           mathFunc("logb", math.Logb),
		"gamma":          mathFunc("gamma", math.Gamma),
		"tgamma":         mathFunc("tgamma", math.Gamma),
		"lgamma":         mathFunc("lgamma", func(v float64) float64 { v, _ = math.Lgamma(v); return v }),
		"erf":            mathFunc("erf", math.Erf),
		"erfc":           mathFunc("erfc", math.Erfc),
		"j0":             mathFunc("j0", math.J0),
		"j1":             mathFunc("j1", math.J1),
		"y0":             mathFunc("y0", math.Y0),
		"y1":             mathFunc("y1", math.Y1),
		"_type_error":    internalfuncTypeError,
	}
}

func noArgFunc(fn func(interface{}) interface{}) function {
	return func(_ *env, f *Func) func(interface{}) interface{} {
		return func(v interface{}) interface{} {
			if len(f.Args) != 0 {
				return &funcNotFoundError{f}
			}
			return fn(v)
		}
	}
}

func mathFunc(name string, f func(x float64) float64) function {
	return noArgFunc(func(v interface{}) interface{} {
		switch v := v.(type) {
		case int:
			return f(float64(v))
		case float64:
			return f(v)
		default:
			return &funcTypeError{name, v}
		}
	})
}

func funcNull(_ interface{}) interface{} {
	return nil
}

func funcTrue(_ interface{}) interface{} {
	return true
}

func funcFalse(_ interface{}) interface{} {
	return false
}

func funcEmpty(_ interface{}) interface{} {
	return struct{}{}
}

func funcLength(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return len(v)
	case map[string]interface{}:
		return len(v)
	case string:
		return len([]rune(v))
	case int:
		if v >= 0 {
			return v
		}
		return -v
	case float64:
		return math.Abs(v)
	case nil:
		return 0
	default:
		return &funcTypeError{"length", v}
	}
}

func funcUtf8ByteLength(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		return len([]byte(v))
	default:
		return &funcTypeError{"utf8bytelength", v}
	}
}

func funcKeys(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		w := make([]interface{}, len(v))
		for i := range v {
			w[i] = i
		}
		return w
	case map[string]interface{}:
		w := make([]string, len(v))
		var i int
		for k := range v {
			w[i] = k
			i++
		}
		sort.Strings(w)
		u := make([]interface{}, len(v))
		for i, x := range w {
			u[i] = x
		}
		return u
	default:
		return &funcTypeError{"keys", v}
	}
}

func funcHas(env *env, f *Func) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		if len(f.Args) != 1 {
			return &funcNotFoundError{f}
		}
		return mapIterator(env.applyPipe(f.Args[0], unitIterator(v)), func(x interface{}) interface{} {
			switch v := v.(type) {
			case []interface{}:
				switch x := x.(type) {
				case int:
					return 0 <= x && x < len(v)
				case float64:
					return 0 <= int(x) && int(x) < len(v)
				default:
					return &hasKeyTypeError{v, x}
				}
			case map[string]interface{}:
				switch x := x.(type) {
				case string:
					_, ok := v[x]
					return ok
				default:
					return &hasKeyTypeError{v, x}
				}
			default:
				return &hasKeyTypeError{v, x}
			}
		})
	}
}

func funcToNumber(v interface{}) interface{} {
	switch v := v.(type) {
	case int, uint, float64:
		return v
	case string:
		var x float64
		if err := json.Unmarshal([]byte(v), &x); err != nil {
			return fmt.Errorf("%s: %q", err, v)
		}
		return x
	default:
		return &funcTypeError{"tonumber", v}
	}
}

func funcType(v interface{}) interface{} {
	return typeof(v)
}

func funcExplode(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		return explode(v)
	default:
		return &funcTypeError{"explode", v}
	}
}

func explode(s string) []interface{} {
	rs := []int32(s)
	xs := make([]interface{}, len(rs))
	for i, r := range rs {
		xs[i] = int(r)
	}
	return xs
}

func funcImplode(v interface{}) interface{} {
	switch v := v.(type) {
	case []interface{}:
		return implode(v)
	default:
		return &funcTypeError{"implode", v}
	}
}

func implode(v []interface{}) interface{} {
	var rs []rune
	for _, r := range v {
		switch r := r.(type) {
		case int:
			rs = append(rs, rune(r))
		case float64:
			rs = append(rs, rune(r))
		default:
			return &funcTypeError{"implode", v}
		}
	}
	return string(rs)
}

func funcJoin(env *env, f *Func) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		if len(f.Args) != 1 {
			return &funcNotFoundError{f}
		}
		return mapIterator(env.applyPipe(f.Args[0], unitIterator(v)), func(x interface{}) interface{} {
			switch v := v.(type) {
			case []interface{}:
				switch x := x.(type) {
				case string:
					var s strings.Builder
					for i, v := range v {
						if i > 0 {
							s.WriteString(x)
						}
						if v != nil {
							if _, ok := v.(float64); ok {
								s.WriteString(funcToJSON(normalizeValues(v)).(string))
							} else {
								s.WriteString(fmt.Sprint(v))
							}
						}
					}
					return s.String()
				default:
					return &funcTypeError{"join", v}
				}
			default:
				return &funcTypeError{"join", v}
			}
		})
	}
}

func funcToJSON(v interface{}) interface{} {
	xs, err := json.Marshal(v)
	if err != nil {
		xs, err = json.Marshal(normalizeValues(v))
		if err != nil {
			return err
		}
	}
	return string(xs)
}

func funcFromJSON(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		var w interface{}
		err := json.Unmarshal([]byte(v), &w)
		if err != nil {
			return err
		}
		return w
	default:
		return &funcTypeError{"fromjson", v}
	}
}

func internalfuncTypeError(env *env, f *Func) func(interface{}) interface{} {
	return func(v interface{}) interface{} {
		if len(f.Args) != 1 {
			return &funcNotFoundError{f}
		}
		return mapIterator(env.applyPipe(f.Args[0], unitIterator(v)), func(x interface{}) interface{} {
			return &funcTypeError{x.(string), v}
		})
	}
}
