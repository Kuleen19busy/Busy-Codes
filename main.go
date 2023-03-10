package main

import (
	"encoding/json"
	"fmt"
	"reflect"
)

func check(value interface{}) error {
	fmt.Println("")
	fmt.Println("Iterating over the key:")
	fmt.Println("")
	fmt.Println("Input Type: ", reflect.TypeOf(value).Kind())

	switch reflect.TypeOf(value).Kind() {
	case reflect.Slice:
		arr := value.([]interface{})
		for _, val := range arr {
			check(val)
		}
	case reflect.Map:
		fmt.Println("")
		fmt.Println("Unpacking Json data:")
		fmt.Println("")
		jsonData, e := json.Marshal(value)
		if e != nil {
			return e
		}
		iterateInputString(jsonData)
	default:
		fmt.Println("value: ", value)
	}
	fmt.Println("")

	return nil
}

func iterateArray(arr []interface{}) error {
	for i, c := range arr {
		fmt.Println("")
		fmt.Println("index", i)
		check(c)
	}
	return nil
}

func iterateInputString(input_data []byte) error {
	input_string := make(map[string]interface{})
	e := json.Unmarshal(input_data, &input_string)
	if e != nil {
		return e
	}

	for key, value := range input_string {
		fmt.Println("key: ", key)

		e = check(value)
		if e != nil {
			return e
		}

	}
	return nil

}

func main() {
	input_string := []byte(`{"name" : "Tolexo Online Pvt. Ltd","age_in_years" : 8.5,"origin" : "Noida","head_office" : "Noida, Uttar Pradesh","address" : [{"street" : "91 Springboard","landmark" : "Axis Bank","city" : "Noida","pincode" : 201301,"state" : "Uttar Pradesh"},{"street" : "91 Springboard","landmark" : "Axis Bank","city" : "Noida","pincode" : 201301,"state" : "Uttar Pradesh"}],"sponsers" : {"name" : "One"},"revenue" : "19.8 million$","no_of_employee" : 630,"str_text" : ["one","two"],"int_text" : [1,3,4]}`)
	var e error
	if e = iterateInputString(input_string); e != nil {
		fmt.Println("invalid input", e)
	}
}
