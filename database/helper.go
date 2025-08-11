package database

import (
	"reflect"

	"github.com/DerKnerd/gorp"
)

func Get[R any](keys ...any) (*R, error) {
	var r R
	res, err := dbMap.Get(&r, keys...)
	if res != nil {
		return res.(*R), err
	}

	return nil, err
}

func SelectOne[R any](query string, keys ...any) (R, error) {
	var r R
	err := dbMap.SelectOne(&r, query, keys...)

	return r, err
}

func Select[R any](query string, keys ...any) ([]R, error) {
	r := new(R)
	res, err := dbMap.Select(&r, query, keys...)
	if err != nil {
		return nil, err
	}

	result := make([]R, len(res))
	for idx, item := range res {
		elem := item.(*R)
		result[idx] = *elem
	}

	return result, nil
}

func toPtr[R any](value any) any {
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr {
		return value
	} else {
		elem := value.(R)
		return &elem
	}
}

func Insert[R any](items ...R) error {
	elems := make([]interface{}, len(items))
	for idx, item := range items {
		elems[idx] = toPtr[R](item)
	}

	return GetDbMap().Insert(elems...)
}

func InsertTx[R any](tx *gorp.Transaction, items ...R) error {
	elems := make([]interface{}, len(items))
	for idx, item := range items {
		elems[idx] = toPtr[R](item)
	}

	return tx.Insert(elems...)
}

func Update[R any](items ...R) (int64, error) {
	elems := make([]interface{}, len(items))
	for idx, item := range items {
		elems[idx] = toPtr[R](item)
	}

	return GetDbMap().Update(elems...)
}

func UpdateTx[R any](tx *gorp.Transaction, items ...R) (int64, error) {
	elems := make([]interface{}, len(items))
	for idx, item := range items {
		elems[idx] = toPtr[R](item)
	}

	return tx.Update(elems...)
}

func Delete[R any](keys ...R) (int64, error) {
	elems := make([]interface{}, len(keys))
	for idx, item := range keys {
		elems[idx] = toPtr[R](item)
	}

	return GetDbMap().Delete(elems...)
}

func DeleteTx[R any](tx *gorp.Transaction, keys ...R) (int64, error) {
	elems := make([]interface{}, len(keys))
	for idx, item := range keys {
		elems[idx] = toPtr[R](item)
	}

	return tx.Delete(elems...)
}
