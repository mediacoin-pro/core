package enc

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// String returns object as string (encode to json)
func String(v interface{}) string {
	switch s := v.(type) {
	case string:
		return s
	case fmt.Stringer:
		return s.String()
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

// JSON returns object as json-string
func JSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func IndentJSON(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

func DecodeJSON(data []byte) (v interface{}) {
	if err := json.Unmarshal(data, &v); err != nil {
		return nil //"<Invalid-JSON>"
	}
	return
}

//func NumToSI(size float64, suffix string ) string {
//
//	switch {
//	case size >= consts.EiB:
//		return fmt.Sprintf("%.2fE", float64(size)/float64(consts.EiB))
//
//	case size >= consts.PiB:
//		return fmt.Sprintf("%.2f PiB", float64(size)/float64(consts.PiB))
//
//	case size >= consts.TiB:
//		return fmt.Sprintf("%.2f TiB", float64(size)/float64(consts.TiB))
//
//	case size >= consts.GiB:
//		return fmt.Sprintf("%.2f GiB", float64(size)/float64(consts.GiB))
//
//	case size >= consts.MiB:
//		return fmt.Sprintf("%.2f MiB", float64(size)/float64(consts.MiB))
//
//	case size >= consts.KiB:
//		return fmt.Sprintf("%.2f KiB", float64(size)/float64(consts.KiB))
//
//	default:
//		return fmt.Sprintf("%d B", size)
//	}
//}

var sizeSfxs = []string{" B", " KB", " MB", " GB", " TB", " PB", " EB"}

func BinarySizeToString(size int64) string {
	f, sfx, pfx := float64(size), 0, ""
	if f < 0 {
		f, pfx = -f, "-"
	}
	for f >= 1000 {
		f /= 1024
		sfx++
	}
	v := strconv.FormatFloat(f, 'f', 3, 64)
	if len(v) > 4 {
		v = v[:4]
	}
	if strings.IndexByte(v, '.') > 0 {
		for v[len(v)-1] == '0' { // trim right '0'
			v = v[:len(v)-1]
		}
	}
	if v[len(v)-1] == '.' { // trim right '.'
		v = v[:len(v)-1]
	}
	return pfx + v + sizeSfxs[sfx]
}
