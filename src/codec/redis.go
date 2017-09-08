package codec

import (
	"reflect"
	"strings"
	"strconv"
	"bytes"
	"errors"
	//"encoding/binary"
	//"log"
)

const (
	END = "\r\n"
)

type RedisError struct {
	err string
}

func (e *RedisError) Error() string {
	return e.err
}

func Serialize(v interface{}) (string){
	t := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	switch t.Kind() {
	case reflect.Bool:
		if v.(bool) {
			return ":1" + END
		}else {
			return ":0" + END
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ":" + strconv.FormatInt(val.Int(), 10) + END
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return ":" + strconv.FormatUint(val.Uint(), 10) + END
	case reflect.Float32, reflect.Float64:
		return "+" + strconv.FormatFloat(val.Float(), 'f', -1, t.Bits()) + END
	case reflect.String:
		return "+" + val.String() + END
	case reflect.Slice, reflect.Array:
		n := val.Len()
		strBuf := make([]string, n + 1)
		strBuf[0] = "*" + strconv.FormatInt(int64(n), 10) + END
		for i:=0;i<n;i++ {
			indexVal := val.Index(i)
			if indexVal.Type().Kind() == reflect.String {
				strBuf[i+1] = "$" + strconv.FormatInt(int64(indexVal.Len()), 10) + END + indexVal.String() + END
			}else{
				strBuf[i+1] = Serialize(indexVal.Interface())
			}
		}
		return strings.Join(strBuf, "")
	default:
		return "$-1" + END
	}
	return "$-1" + END//null
}

func Unserialize(buf []byte) (interface{}, []byte, error) {
	strLen := len(buf)
	if strLen < 3 {
		return nil, nil, errors.New("not enough sequences to unserialize, len: " + strconv.FormatInt(int64(strLen), 10))
	}
	switch buf[0] {
	case ':':
		ret := bytes.SplitN(bytes.TrimSpace(buf[1:]), []byte("\r\n"), 2)
		num, err := strconv.ParseInt(string(ret[0]), 10, 64)
		if err != nil {
			return nil, nil, err
		}
		if len(ret) > 1 {
			return num, ret[1], nil
		}
		return num, nil, nil
	case '+':
		ret := bytes.SplitN(bytes.TrimSpace(buf[1:]), []byte("\r\n"), 2)
		return string(ret[0]), ret[1], nil
	case '-':
		ret := bytes.SplitN(bytes.TrimSpace(buf[1:]), []byte("\r\n"), 2)
		if len(ret) > 1 {
			return nil, ret[1], errors.New(string(ret[0]))
		}
		return errors.New(string(ret[0])), nil, nil
	case '$':
		if string(buf[1:5]) == "-1\r\n" {
			return nil, buf[5:], nil
		}
		ret := bytes.SplitN(bytes.TrimSpace(buf[1:]), []byte("\r\n"), 2)
		num, err:= strconv.Atoi(string(ret[0]))
		if err != nil {
			return nil, nil, err
		}
		return string(ret[1][0:num]), ret[1][num:], nil
	case '*':
		ret := bytes.SplitN(bytes.TrimSpace(buf[1:]), []byte("\r\n"), 2)
		num, err := strconv.Atoi(string(ret[0]))
		if err != nil {
			return nil, buf, err
		}
		if num == -1 {
			return nil, ret[1], nil
		}
		arr := make([]interface{}, num)
		tmp := ret[1]
		for i := 0; i < int(num); i++ {
			if len(ret) == 1 {
				break
			}
			val, tmpBuf, err := Unserialize(tmp)
			if err != nil {
				return nil, buf, err
			}
			tmp = bytes.TrimSpace(tmpBuf)
			arr[i] = val
		}
		return arr, tmp, nil
	}
	return nil, nil, nil
}