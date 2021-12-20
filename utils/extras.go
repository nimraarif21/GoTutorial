package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gotutorial/api/models"
)



var TaskFields = getTaskFields()

func getTaskFields() []string {
    var field []string

    v := reflect.ValueOf(models.Task{})
    for i := 0; i < v.Type().NumField(); i++ {
        field = append(field, v.Type().Field(i).Tag.Get("json"))
    }
    return field
}


func ValidateAndReturnSortQuery(sortBy string) (string, error) {
    splits := strings.Split(sortBy, ".")
    if len(splits) != 2 {
        return "", errors.New("malformed sortBy query parameter, should be field.orderdirection")
    }

    field, order := splits[0], splits[1]

    if order != "desc" && order != "asc" {
        return "", errors.New("malformed orderdirection in sortBy query parameter, should be asc or desc")
    }

    if !stringInSlice(TaskFields, field) {
        return "", errors.New("unknown field in sortBy query parameter")
    }

    return fmt.Sprintf("%s %s", field, strings.ToUpper(order)), nil

}

func stringInSlice(strSlice []string, s string) bool {
    for _, v := range strSlice {
        if v == s {
            return true
        }
    }

    return false
}

func ValidateAndReturnFilterMap(filter string) (map[string]string, error) {
	splits := strings.Split(filter, ".")
	if len(splits) != 2 {
			return nil, errors.New("malformed sortBy query parameter, should be field.orderdirection")
	}
	field, value := splits[0], splits[1]
	if !stringInSlice(TaskFields, field) {
			return nil, errors.New("unknown field in filter query parameter")
	}
	return map[string]string{field: value}, nil
}