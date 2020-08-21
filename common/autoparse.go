package common

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"reflect"
	"strconv"
)

func AutoParse(c interface{}) {

	c_value := reflect.ValueOf(c)

	c_type := c_value.Elem().Type()

	for i := 0; i < c_type.NumField(); i++ {

		field := c_type.Field(i)
		fg := field.Tag.Get("flag")
		fd := field.Tag.Get("default")
		fu := field.Tag.Get("usage")

		switch field.Type.Kind() {
		case reflect.String:
			ontology := c_value.Elem().Field(i).Addr().Interface().(*string)
			flag.StringVar(ontology, fg, fd, fu)

		case reflect.Int:
			ontology := c_value.Elem().Field(i).Addr().Interface().(*int)
			fd_int, err := strconv.Atoi(fd)
			if err != nil {
				panic(err)
			}
			flag.IntVar(ontology, fg, fd_int, fu)
		case reflect.Int64:
			ontology := c_value.Elem().Field(i).Addr().Interface().(*int64)
			fd_int, err := strconv.Atoi(fd)
			if err != nil {
				panic(err)
			}
			flag.Int64Var(ontology, fg, int64(fd_int), fu)

		case reflect.Bool:
			ontology := c_value.Elem().Field(i).Addr().Interface().(*bool)
			flag.BoolVar(ontology, fg, fd == "true", fu)

		case reflect.Float64:
			ontology := c_value.Elem().Field(i).Addr().Interface().(*float64)
			fd_float, err := strconv.ParseFloat(fd, 64)
			if err != nil {
				panic(err)
			}
			flag.Float64Var(ontology, fg, fd_float, fu)

		default:
			panic(field.Name + " this field type is not support")
		}

	}

	flag.Parse()

}

func Load(path string, conf interface{}) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		panic(err)
	}

}
