package eventutil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func UnmarshalKVMap(kvs map[string]string, v interface{}) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New(fmt.Sprintf("wrong value to unmarshal %v", reflect.TypeOf(v)))
	}

	rv = rv.Elem()

	for i := 0; i < rv.NumField(); i++ {
		typInfo := rv.Type().Field(i)
		f := rv.Field(i)
		tag := typInfo.Tag.Get("json")
		if tag == "" {
			tag = strings.ToLower(typInfo.Name)
		}

		value, ok := kvs[tag]
		if ok {
			if f.Kind() == reflect.Slice && f.Type() != reflect.TypeOf(types.AccAddress{}) {
				// use `,` to be seq
				strs := strings.Split(value, ",")
				for _, str := range strs {
					str = strings.TrimSpace(str)
					elem := reflect.New(f.Type().Elem()).Elem()
					if err := populate(elem, str); err != nil {
						return fmt.Errorf("%s: %v", tag, err)
					}
					f.Set(reflect.Append(f, elem))
				}
			} else if err := populate(f, value); err != nil {
				return fmt.Errorf("%s: %v", tag, err)
			}
		}
	}

	return nil
}

func UnmarshalEvent(evt sdk.Event, v interface{}) error {
	kvs := make(map[string]string, len(evt.Attributes))
	for _, attr := range evt.Attributes {
		if _, ok := kvs[string(attr.Key)]; ok {
			panic(errors.New(fmt.Sprintf("event type has two equal key %s", string(attr.Key))))
		}
		kvs[string(attr.Key)] = string(attr.Value)
	}

	return UnmarshalKVMap(kvs, v)
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)
		return nil

	case reflect.Int:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
		return nil

	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
		return nil

	case reflect.Uint:
		u, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(u)
		return nil

	case reflect.Uint64:
		u, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(u)
		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)
		return nil
	default:
		return populateTypes(v, value)
	}
}

func populateTypes(v reflect.Value, value string) error {
	switch v.Type() {
	case reflect.TypeOf(types.Name{}):
		n, err := types.NewName(value)
		if err != nil {
			return errors.Wrapf(err, "populate %s from name", value)
		}
		v.Set(reflect.ValueOf(n))

	case reflect.TypeOf(types.AccAddress{}):
		accAddr, err := types.AccAddressFromBech32(value)
		if err != nil {
			return errors.Wrapf(err, "populate %s from accaddress", value)
		}
		v.Set(reflect.ValueOf(accAddr))

	case reflect.TypeOf(types.AccountID{}):
		id, err := types.NewAccountIDFromStr(value)
		if err != nil {
			return errors.Wrapf(err, "populate %s from account id", value)
		}
		v.Set(reflect.ValueOf(id))

	case reflect.TypeOf(types.Coin{}):
		coin, err := types.ParseCoin(value)
		if err != nil {
			return errors.Wrapf(err, "populate %s coin", value)
		}
		v.Set(reflect.ValueOf(coin))

	default:
		return fmt.Errorf("unsupported kind %s", v.Type())
	}

	return nil
}
