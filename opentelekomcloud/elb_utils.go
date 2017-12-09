package opentelekomcloud

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform/helper/schema"
)

func chooseCESClient(d *schema.ResourceData, config *Config) (*gophercloud.ServiceClient, error) {
	return config.loadCESClient(GetRegion(d, config))
}

func isCESResourceNotFound(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(gophercloud.ErrDefault404)
	return ok
}

func buildCESCreateParam(opts interface{}, d *schema.ResourceData) error {
	return buildCESCUParam(opts, d, false)
}

func buildCESUpdateParam(opts interface{}, d *schema.ResourceData) error {
	return buildCESCUParam(opts, d, true)
}

func buildCESCUParam(opts interface{}, d *schema.ResourceData, buildUpdate bool) error {
	optsValue := reflect.ValueOf(opts)
	if optsValue.Kind() != reflect.Ptr {
		return fmt.Errorf("parameter of opts should be a pointer")
	}
	optsValue = optsValue.Elem()
	if optsValue.Kind() != reflect.Struct {
		return fmt.Errorf("parameter must be a pointer to a struct")
	}

	optsType := reflect.TypeOf(opts)
	optsType = optsType.Elem()

	value := make(map[string]interface{})
	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			return fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		param := strings.Split(tag, ",")[0]
		if buildUpdate && !d.HasChange(param) {
			continue
		}

		if d.Get(param) == nil {
			log.Printf("[DEBUG] param:%s is not set", param)
			continue
		}

		value[param] = d.Get(param)
		/*
			switch v.Kind() {
			case reflect.String:
				v.SetString(d.Get(param).(string))
			case reflect.Int:
				v.SetInt(int64(d.Get(param).(int)))
			case reflect.Bool:
				v.SetBool(d.Get(param).(bool))
			case reflect.Slice:
				s := d.Get(param).([]interface{})

				switch v.Type().Elem().Kind() {
				case reflect.String:
					t := make([]string, len(s))
					for i, iv := range s {
						t[i] = iv.(string)
					}
					v.Set(reflect.ValueOf(t))
				default:
					return fmt.Errorf("unknown type of item %v: %v", v, v.Type().Elem().Kind())
				}
			case reflect.Struct:
				e := buildStruct(&v, f.Type, d.Get(param).(map[string]interface{}))
				if e != nil {
					return e
				}
			default:
				return fmt.Errorf("unknown type of item %v: %v", v, v.Kind())
			}
		*/
	}
	return buildStruct(&optsValue, optsType, value)
}

func buildStruct(optsValue *reflect.Value, optsType reflect.Type, value map[string]interface{}) error {
	log.Printf("[DEBUG] buildStruct:: optsValue=%v, optsType=%v, value=%#v\n", optsValue, optsType, value)

	for i := 0; i < optsValue.NumField(); i++ {
		v := optsValue.Field(i)
		f := optsType.Field(i)
		tag := f.Tag.Get("json")
		if tag == "" {
			return fmt.Errorf("can not convert for item %v: without of json tag", v)
		}
		param := strings.Split(tag, ",")[0]
		log.Printf("[DEBUG] buildStruct:: convert for param:%s", param)
		if _, e := value[param]; !e {
			log.Printf("[DEBUG] param:%s was not supplied", param)
			continue
		}

		switch v.Kind() {
		case reflect.String:
			v.SetString(value[param].(string))
		case reflect.Int:
			v.SetInt(int64(value[param].(int)))
		case reflect.Int64:
			v.SetInt(value[param].(int64))
		case reflect.Bool:
			v.SetBool(value[param].(bool))
		case reflect.Slice:
			s := value[param].([]interface{})

			switch v.Type().Elem().Kind() {
			case reflect.String:
				t := make([]string, len(s))
				for i, iv := range s {
					t[i] = iv.(string)
				}
				v.Set(reflect.ValueOf(t))
			case reflect.Struct:
				t := reflect.MakeSlice(f.Type, len(s), len(s))
				for i, iv := range s {
					rv := t.Index(i)
					e := buildStruct(&rv, f.Type.Elem(), iv.(map[string]interface{}))
					if e != nil {
						return e
					}
				}
				v.Set(t)

			default:
				return fmt.Errorf("unknown type of item %v: %v", v, v.Type().Elem().Kind())
			}
		case reflect.Struct:
			log.Printf("[DEBUG] buildStruct:: convert struct for param %s: %#v", param, value[param])
			//The corresponding type of Struct is TypeList in Terrafrom
			var p map[string]interface{}

			v0, ok := value[param].([]interface{})
			if ok {
				p, ok = v0[0].(map[string]interface{})
			} else {
				p, ok = value[param].(map[string]interface{})
			}
			if !ok {
				return fmt.Errorf("can not convert to (map[string]interface{}) for param %s: %#v", param, value[param])
			}

			e := buildStruct(&v, f.Type, p)
			if e != nil {
				return e
			}

		default:
			return fmt.Errorf("unknown type of item %v: %v", v, v.Kind())
		}
	}
	return nil
}
