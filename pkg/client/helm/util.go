package helm

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/release"
)

func convertArgsToStruct(args map[string]interface{}, obj interface{}) error {
	v := reflect.ValueOf(obj)

	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("non-pointer %v", v.Type().String())
	}
	if v.IsNil() {
		return fmt.Errorf("nil %v", v.Type().String())
	}

	v = v.Elem()

	for i := 0; i < v.NumField(); i++ {
		var tagName string
		var isOptional bool
		f := v.Type().Field(i)
		if paramTag := f.Tag.Get("param"); paramTag != "" {
			tags := strings.Split(paramTag, ",")
			tagName = tags[0]
			if len(tags) > 1 && tags[1] == "optional" {
				isOptional = true
			}
		} else {
			tagName = camelCaseToSnakeCase(f.Name)
		}

		if value, ok := args[tagName]; !ok {
			if isOptional {
				continue
			} else {
				return fmt.Errorf("argument(%s) does not exists", tagName)
			}
		} else {
			argValue := reflect.ValueOf(value)
			if v.Field(i).Type() == argValue.Type() {
				if v.Field(i).CanSet() {
					v.Field(i).Set(argValue)
				}
			} else {
				if argValue.Kind() == reflect.Float64 && v.Field(i).Kind() == reflect.Int {
					if v.Field(i).CanSet() && argValue.CanConvert(v.Field(i).Type()) {
						v.Field(i).Set(argValue.Convert(v.Field(i).Type()))
					}
				} else {
					return fmt.Errorf("argument(%s)'s type is invalid: want(%v), got(%v)", tagName, v.Field(i).Type(), argValue.Type())
				}
			}
		}
	}
	return nil
}

func camelCaseToSnakeCase(input string) string {
	var res string
	for k, v := range input {
		if k == 0 {
			res = strings.ToLower(string(v))
		} else {
			if v >= 'A' && v <= 'Z' {
				res += "_" + strings.ToLower(string(v))
			} else {
				res += string(v)
			}
		}
	}

	return res
}

func extractResultFrom(rel *release.Release) (map[string]interface{}, error) {
	if rel == nil {
		return nil, fmt.Errorf("release.Release is nil")
	}

	res := make(map[string]interface{})

	res["name"] = rel.Name
	res["namespace"] = rel.Namespace
	res["status"] = rel.Info.Status.String()
	res["revision"] = rel.Version

	if rel.Info != nil {
		if !rel.Info.LastDeployed.IsZero() {
			res["last_deployed"] = rel.Info.LastDeployed.Format(time.ANSIC)
		}
		if len(rel.Info.Notes) > 0 {
			res["notes"] = strings.TrimSpace(rel.Info.Notes)
		}
	}

	if len(rel.Config) > 0 {
		res["user_supplied_values"] = rel.Config
	}

	return res, nil
}

func extractHistoryResultFrom(rels []*release.Release) ([]map[string]interface{}, error) {
	if rels == nil {
		return nil, fmt.Errorf("release.Release is nil")
	}

	var results []map[string]interface{}

	for _, rel := range rels {
		res := make(map[string]interface{})

		res["namespace"] = rel.Namespace
		res["release_name"] = rel.Name
		res["revision"] = rel.Version

		if rel.Chart != nil {
			res["app_version"] = rel.Chart.AppVersion()
			res["chart_name"] = rel.Chart.Name()

			if rel.Chart.Metadata != nil {
				res["chart_version"] = rel.Chart.Metadata.Version
			} else {
				res["chart_version"] = "MISSING"
			}
		} else {
			res["app_version"] = "MISSING"
			res["chart_name"] = "MISSING"
		}

		if rel.Info != nil {
			res["status"] = rel.Info.Status.String()
			res["description"] = rel.Info.Description

			if !rel.Info.LastDeployed.IsZero() {
				res["updated"] = rel.Info.LastDeployed
			}
		}

		results = append(results, res)
	}

	return results, nil
}
