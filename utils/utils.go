package utils

import (
	"errors"
	"strconv"
)

type PathInfo struct {
    Id int
}

func GetPathValues(ps []string) (PathInfo, error){
    r := PathInfo{
	Id: 0,
    }

    if len(ps) > 3 {
	if ps[3] != "" {
	    err := errors.New("not found")
	    return r, err
	}
    }

    id, err := strconv.Atoi(ps[2])
    if err != nil {
	err := errors.New("not an integer")
	return r, err
    }
    r.Id = id

    return r, err
}
