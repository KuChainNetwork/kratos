package types

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// testNameConvert test the name convert from string
func testNameConvert(t *testing.T, str string) {
	Convey(fmt.Sprintf("test name conv %s", str), t, func() {
		name, err := NewName(str)
		So(err, ShouldEqual, nil)

		//debugShowName(t, name)

		nameStr := name.String()
		So(nameStr, ShouldEqual, str)

		So(name.Value[0], ShouldEqual, 1)
		So(name.Value[1], ShouldEqual, 1)
		So(name.Value[2], ShouldEqual, len(str))
	})
}

func TestName_String(t *testing.T) {
	testNameConvert(t, "k")
	testNameConvert(t, "ku")
	testNameConvert(t, "kuc")
	testNameConvert(t, "kuch")
	testNameConvert(t, "kuchain")
	testNameConvert(t, "kuchain@123")
	testNameConvert(t, "kuchain1122334")
	testNameConvert(t, "kuchainvcdf2322a3")
	testNameConvert(t, ".....sdsdsd......")
	testNameConvert(t, "_________________")
	testNameConvert(t, "_")
	testNameConvert(t, "3df@...____")
}

func TestName_ParseErr(t *testing.T) {
	Convey("test string valid", t, func() {
		_, err := NewName("kuchain11111111111")
		So(err, ShouldEqual, ErrNameStrNoValid)

		_, err1 := NewName("@")
		So(err1, ShouldEqual, ErrNameStrNoValid)

		_, err2 := NewName("@12221")
		So(err2, ShouldEqual, ErrNameStrNoValid)

		_, err3 := NewName("sadfdsd@")
		So(err3, ShouldEqual, ErrNameStrNoValid)

		_, err4 := NewName("@asdf@")
		So(err4, ShouldEqual, ErrNameStrNoValid)

		_, err5 := NewName("asdfsdf@fsdf@")
		So(err5, ShouldEqual, ErrNameStrNoValid)

		_, err6 := NewName("Dfdsdsd")
		So(err6, ShouldEqual, ErrNameStrNoValid)

		_, err7 := NewName("!dsfsddsdff")
		So(err7, ShouldEqual, ErrNameStrNoValid)

		_, err8 := NewName("-dsfsddsdff")
		So(err8, ShouldEqual, ErrNameStrNoValid)
	})
}

func TestName_Char2Byte(t *testing.T) {
	Convey("char 2 bytes", t, func() {
		for _, r := range "abcdefghijklmnopqrstuvwxyz" {
			c := byte(r)
			So(byte2char(char2byte(c)), ShouldEqual, c)
		}

		for _, r := range "0123456789" {
			c := byte(r)
			So(byte2char(char2byte(c)), ShouldEqual, c)
		}

		for _, r := range "@._" {
			c := byte(r)
			So(byte2char(char2byte(c)), ShouldEqual, c)
		}

		So(char2byte(byte('a')), ShouldEqual, 1)
		So(char2byte(byte('z')), ShouldEqual, 26)
		So(char2byte(byte('0')), ShouldEqual, 32)
		So(char2byte(byte('9')), ShouldEqual, 41)
		So(char2byte(byte('.')), ShouldEqual, CharValueDot)
		So(char2byte(byte('_')), ShouldEqual, CharValueUnderline)
		So(char2byte(byte('@')), ShouldEqual, CharValueSeq)
	})

	Convey("char error", t, func() {
		chars := "@abcdefghijklmnopqrstuvwxyz*****0123456789*******._"
		for v := 0; v < 64; v++ {
			if v < len(chars) {
				So(byte2char(byte(v)), ShouldEqual, chars[v])
			} else {
				So(byte2char(byte(v)), ShouldEqual, '*')
			}
		}
	})
}

func BenchmarkNameParsePerformance(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < 1000000; i++ {
		MustName("kucha23sds3s212a3")
	}
}

func testNameNoEq(t *testing.T, l, r string) {
	Convey(fmt.Sprintf("test name no eq %s %s", l, r), t, func() {
		ln := MustName(l)
		rn := MustName(r)
		So(ln.Eq(rn), ShouldEqual, false)
	})

	Convey(fmt.Sprintf("test name eq %s %s", l, r), t, func() {
		ln := MustName(l)
		So(ln.Eq(ln), ShouldEqual, true)

		rn := MustName(r)
		So(rn.Eq(rn), ShouldEqual, true)
	})
}

func TestNameEq(t *testing.T) {
	testNameNoEq(t, "@@", "@")
	testNameNoEq(t, "kuchain", "kuchai")
	testNameNoEq(t, "kuchain1212121212", "kuchain1212121213")
	testNameNoEq(t, "kuchain1212121212", "kuchain121212121@")
	testNameNoEq(t, "@@@@@@@@@@@@@@@@@", "1@@@@@@@@@@@@@@@@")
}

func TestName_Foreach(t *testing.T) {
	nameStr := "kuchainvcdf2322a3"
	Convey("test for each func", t, func() {
		name := MustName(nameStr)
		foreachStr := make([]byte, 0, NameStrLenMax+1)
		name.Foreach(func(c byte) bool {
			foreachStr = append(foreachStr, c)
			return true
		})
		So(string(foreachStr), ShouldEqual, nameStr)
	})
}

func TestName_Len(t *testing.T) {
	nameStr1 := "kuchainvcdf2322a3"
	nameStr2 := "k"

	Convey("test name len", t, func() {
		So(len(nameStr1), ShouldEqual, MustName(nameStr1).Len())
		So(len(nameStr2), ShouldEqual, MustName(nameStr2).Len())
	})
}

func TestNewNameFromBytes(t *testing.T) {
	nameStr := "kuchainvcdf2322a3"
	Convey("test new byte", t, func() {
		n := MustName(nameStr)
		nn := NewNameFromBytes(n.Value[:])
		So(n.Eq(nn), ShouldEqual, true)
	})
}
