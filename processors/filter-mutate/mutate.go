//go:generate bitfanDoc
// mutate filter allows to perform general mutations on fields. You can rename, remove, replace, and modify fields in your event.
package mutate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/clbanning/mxj"
	"bitfan/processors"
)

const (
	PORT_SUCCESS = 0
)

// Performs mutations on fields
func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base
	opt *options
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	Add_tag []string

	// Convert a field’s value to a different type, like turning a string to an integer.
	// If the field value is an array, all members will be converted. If the field is a hash,
	// no action will be taken.
	// If the conversion type is boolean, the acceptable values are:
	// True: true, t, yes, y, and 1
	// False: false, f, no, n, and 0
	// If a value other than these is provided, it will pass straight through and log a warning message.
	// Valid conversion targets are: integer, float, string, and boolean.
	Convert map[string]string

	// Convert a string field by applying a regular expression and a replacement. If the field is not a string, no action will be taken.
	// This configuration takes an array consisting of 3 elements per field/substitution.
	// Be aware of escaping any backslash in the config file.
	Gsub []string

	// Join an array with a separator character. Does nothing on non-array fields
	Join map[string]string

	// Convert a value to its lowercase equivalent
	Lowercase []string

	// Merge two fields of arrays or hashes. String fields will be automatically be converted into an array
	Merge map[string]string

	// If this filter is successful, remove arbitrary fields from this event.
	Remove_field []string

	// If this filter is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	Remove_tag []string

	// Rename key on one or more fields
	Rename map[string]string

	// Replace a field with a new value. The new value can include %{foo} strings to
	// help you build a new value from other parts of the event
	Replace map[string]interface{}

	// Split a field to an array using a separator character. Only works on string fields
	Split map[string]string

	// Strip whitespace from processors. NOTE: this only works on leading and trailing whitespace
	Strip []string

	// Update an existing field with a new value. If the field does not exist, then no action will be taken
	Update map[string]interface{}

	// Convert a value to its uppercase equivalent
	Uppercase []string

	// remove all fields, except theses fields (work only with first level fields)
	Remove_all_but []string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Receive(e processors.IPacket) error {
	processors.AddFields(p.opt.Add_field, e.Fields())
	processors.AddTags(p.opt.Add_tag, e.Fields())
	UpdateFields(p.opt.Update, e.Fields())
	UpdateFields(p.opt.Replace, e.Fields())
	processors.RemoveFields(p.opt.Remove_field, e.Fields())
	RenameFields(p.opt.Rename, e.Fields())
	UpperCaseFields(p.opt.Uppercase, e.Fields())
	LowerCaseFields(p.opt.Lowercase, e.Fields())
	RemoveAllButFields(p.opt.Remove_all_but, e.Fields())
	Convert(p.opt.Convert, e.Fields())
	Join(p.opt.Join, e.Fields())
	processors.RemoveTags(p.opt.Remove_tag, e.Fields())
	Gsub(p.opt.Gsub, e.Fields())
	Split(p.opt.Split, e.Fields())
	Strip(p.opt.Strip, e.Fields())
	Merge(p.opt.Merge, e.Fields())

	p.Send(e, PORT_SUCCESS)

	return nil
}

func RemoveAllButFields(fields []string, data *mxj.Map) {
	if len(fields) > 0 {
		cp := mxj.New()
		for _, k := range fields {
			processors.Dynamic(&k, data)
			if value, err := data.ValueForPath(k); err == nil {
				cp.SetValueForPath(value, k)
			}
		}
		*data = cp
	}
}

func UpdateFields(fields map[string]interface{}, data *mxj.Map) {
	for k, v := range fields {
		processors.Dynamic(&k, data)
		if data.Exists(k) {
			switch t := v.(type) {
			case string:
				processors.Dynamic(&t, data)
				v = t
			}
			data.SetValueForPath(v, k)
		}
	}
}

func RenameFields(fields map[string]string, data *mxj.Map) {
	for k, v := range fields {
		processors.Dynamic(&k, data)
		if data.Exists(k) {
			processors.Dynamic(&v, data)
			data.RenameKey(k, v)
		}
	}
}

func UpperCaseFields(fields []string, data *mxj.Map) {
	for _, k := range fields {
		processors.Dynamic(&k, data)
		if value, err := data.ValueForPathString(k); err == nil {
			data.SetValueForPath(strings.ToUpper(value), k)
		}
	}
}

func LowerCaseFields(fields []string, data *mxj.Map) {
	for _, k := range fields {
		processors.Dynamic(&k, data)
		if value, err := data.ValueForPathString(k); err == nil {
			data.SetValueForPath(strings.ToLower(value), k)
		}
	}
}

func Join(fields map[string]string, data *mxj.Map) {
	for path, glue := range fields {
		if !data.Exists(path) {
			continue
		}
		values, _ := data.ValuesForPath(path)
		a := []string{}
		for _, v := range values {
			switch vu := v.(type) {
			case []string:
				as := []string{}
				for _, s := range vu {
					as = append(as, s)
				}
				newValue := strings.Join(as, glue)
				data.SetValueForPath(newValue, path)
			case []int:
				as := []string{}
				for _, i := range vu {
					s := strconv.Itoa(int(i))
					as = append(as, s)
				}
				newValue := strings.Join(as, glue)
				data.SetValueForPath(newValue, path)
			case string:
				a = append(a, vu)
			case int:
				s := strconv.Itoa(int(vu))
				a = append(a, s)
			}
		}

		if len(a) > 0 {
			newValue := strings.Join(a, glue)
			data.SetValueForPath(newValue, path)
		}
	}
}

func Split(fields map[string]string, data *mxj.Map) {
	for path, separator := range fields {
		if !data.Exists(path) {
			continue
		}
		value := data.ValueOrEmptyForPathString(path)

		newValue := strings.Split(value, separator)

		s := make([]interface{}, len(newValue))
		for i, v := range newValue {
			s[i] = v
		}

		data.SetValueForPath(s, path)
	}
}

func Strip(fields []string, data *mxj.Map) {
	for _, path := range fields {
		if value, err := data.ValueForPathString(path); err == nil {
			newValue := strings.TrimSpace(value)
			data.SetValueForPath(newValue, path)
		}

	}
}

func Gsub(fields []string, data *mxj.Map) {
	for i := 0; i < len(fields); i++ {
		fieldname := fields[i]
		i++
		pattern := fields[i]
		i++
		replacement := fields[i]

		if value, err := data.ValueForPathString(fieldname); err == nil {
			if r, err := regexp.Compile(pattern); err == nil {
				newValue := r.ReplaceAllString(value, replacement)
				data.SetValueForPath(newValue, fieldname)
			}
		}
	}
}

func convertValue(value interface{}, kind string) (interface{}, error) {
	switch value.(type) {
	case string:
		switch kind {
		case "integer":
			newValue, err := strconv.Atoi(value.(string))
			if err != nil {
				return nil, err
			}
			return newValue, nil
		case "float":
			newValue, err := strconv.ParseFloat(value.(string), 64)
			if err != nil {
				return nil, err
			}
			return newValue, nil
		case "boolean":
			value = strings.ToLower(value.(string))
			for _, b := range []string{"true", "t", "yes", "y", "1"} {
				if b == value {
					return true, nil
				}
			}
			for _, b := range []string{"false", "f", "no", "n", "0"} {
				if b == value {
					return false, nil
				}
			}
			return nil, fmt.Errorf("unknow boolean %v", value)
		}
	case int:
		switch kind {
		case "string":
			newValue := fmt.Sprintf("%d", value.(int))
			return newValue, nil
		case "float":
			newValue := float64(value.(int))
			return newValue, nil
		case "boolean":
			newValue := false
			if value.(int) > 0 {
				newValue = true
			}
			return newValue, nil
		}
	case float64:
		switch kind {
		case "string":
			newValue := fmt.Sprintf("%f", value.(float64))
			return newValue, nil
		case "integer":
			newValue := int(value.(float64))
			return newValue, nil
		case "boolean":
			newValue := false
			if value.(float64) > 0 {
				newValue = true
			}
			return newValue, nil
		}
	case bool:
		switch kind {
		case "string":
			newValue := "false"
			if value.(bool) == true {
				newValue = "true"
			}
			return newValue, nil
		case "integer":
			var newValue int
			newValue = 0
			if value.(bool) == true {
				newValue = 1
			}
			return newValue, nil
		case "float":
			var newValue float64
			newValue = 0
			if value.(bool) == true {
				newValue = 1
			}
			return newValue, nil
		}
	}
	return nil, fmt.Errorf("unconvertible to %s, value= %v", kind, value)
}

func Convert(fields map[string]string, data *mxj.Map) {
	for path, kind := range fields {
		path = processors.NormalizeNestedPath(path)
		if !data.Exists(path) {
			continue
		}

		values, err := data.ValuesForPath(path)
		if err != nil {
			continue
		}

		if len(values) > 1 {
			newValue := []interface{}{}
			for _, value := range values {
				nValue, err := convertValue(value, kind)
				if err != nil {
					nValue = value
				}
				newValue = append(newValue, nValue)
			}
			data.SetValueForPath(newValue, path)
			continue
		}

		if len(values) == 1 {
			switch v := values[0].(type) {
			case []string:
				newValue := []interface{}{}
				for _, value := range v {
					nValue, err := convertValue(value, kind)
					if err != nil {
						nValue = value
					}
					newValue = append(newValue, nValue)
				}
				data.SetValueForPath(newValue, path)
			case interface{}:
				newValue, err := convertValue(v, kind)
				if err != nil {
					newValue = v
				}
				data.SetValueForPath(newValue, path)
			}

			continue
		}

	}

}

func Merge(fields map[string]string, data *mxj.Map) {
	for path_dst, path_src := range fields {
		if !data.Exists(path_dst) || !data.Exists(path_src) {
			continue
		}

		a := []interface{}{}
		values_src, _ := data.ValuesForPath(path_src)
		for _, value_src := range values_src {
			switch v := value_src.(type) {
			case []string:
				s := make([]interface{}, len(v))
				for i, v1 := range v {
					s[i] = v1
				}
				a = append(a, s...)
			case []int:
				s := make([]interface{}, len(v))
				for i, v1 := range v {
					s[i] = v1
				}
				a = append(a, s...)
			case interface{}:
				a = append(a, v)
			}
		}

		b := []interface{}{}
		values_dst, _ := data.ValuesForPath(path_dst)
		for _, value_dst := range values_dst {
			switch v := value_dst.(type) {
			case []string:
				vi := make([]interface{}, len(v))
				for i, v1 := range v {
					vi[i] = v1
				}
				b = append(b, vi...)
			case []int:
				vi := make([]interface{}, len(v))
				for i, v1 := range v {
					vi[i] = v1
				}
				b = append(b, vi...)
			case interface{}:
				b = append(b, v)
			}
		}

		// Add a to b
		b = append(b, a...)

		// Remove duplicates
		result := []interface{}{}
		seen := map[interface{}]bool{}
		for _, val := range b {
			if _, ok := seen[val]; !ok {
				result = append(result, val)
				seen[val] = true
			}
		}

		data.SetValueForPath(result, path_dst)
	}
}
