package xio

import (
	"bufio"
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"io"
	"strings"
	"testing"
)

func TestOffsetReader_Read(t *testing.T) {
	Convey("Create a new offset reader", t, func() {
		r := &OffsetReader{
			Reader: strings.NewReader("0123456789"),
			Skip:   0,
			Stop:   0,
		}

		Convey("Read a byte", func() {
			p := make([]byte, 1)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 1)
			So(p, ShouldResemble, []byte("0"))
			So(r.Current(), ShouldEqual, 1)
		})

		Convey("Read all bytes", func() {
			p := make([]byte, 10)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 10)
			So(p, ShouldResemble, []byte("0123456789"))
			So(r.Current(), ShouldEqual, 10)

			Convey("Read 1 byte when all bytes already read", func() {
				p = make([]byte, 1)
				_, err := r.Read(p)
				So(err, ShouldEqual, io.EOF)
				So(r.Current(), ShouldEqual, 10)
			})
		})

		Convey("Read All", func() {
			n, err := io.ReadAll(r)
			So(err, ShouldBeNil)
			So(n, ShouldResemble, []byte("0123456789"))
			So(r.Current(), ShouldEqual, 10)
		})
	})

	Convey("Create a new offset reader with skip", t, func() {
		r := &OffsetReader{
			Reader: strings.NewReader("0123456789"),
			Skip:   5,
			Stop:   0,
		}

		Convey("Read a byte", func() {
			p := make([]byte, 1)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 1)
			So(p, ShouldResemble, []byte("5"))
			So(r.Current(), ShouldEqual, 6)
		})

		Convey("Read all bytes", func() {
			p := make([]byte, 5)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 5)
			So(p, ShouldResemble, []byte("56789"))
			So(r.Current(), ShouldEqual, 10)

			Convey("Read 1 byte when all bytes already read", func() {
				p = make([]byte, 1)
				_, err := r.Read(p)
				So(err, ShouldEqual, io.EOF)
				So(r.Current(), ShouldEqual, 10)
			})
		})

	})

	Convey("Create a new offset reader with stop", t, func() {
		r := &OffsetReader{
			Reader: strings.NewReader("0123456789"),
			Skip:   0,
			Stop:   5,
		}

		Convey("Read a byte", func() {
			p := make([]byte, 1)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 1)
			So(p, ShouldResemble, []byte("0"))
			So(r.Current(), ShouldEqual, 1)
		})

		Convey("Read all bytes", func() {
			p := make([]byte, 5)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 5)
			So(p, ShouldResemble, []byte("01234"))
			So(r.Current(), ShouldEqual, 5)

			Convey("Read 1 byte when all bytes already read", func() {
				p := make([]byte, 1)
				_, err := r.Read(p)
				So(err, ShouldEqual, io.EOF)
				So(r.Current(), ShouldEqual, 5)
			})
		})

		Convey("Read 6 bytes", func() {
			p := make([]byte, 6)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 5)
			expected := make([]byte, 6)
			expected[0] = '0'
			expected[1] = '1'
			expected[2] = '2'
			expected[3] = '3'
			expected[4] = '4'
			So(p, ShouldResemble, expected)
			So(r.Current(), ShouldEqual, 5)
		})

	})

	Convey("Create a new offset reader with both skip and stop", t, func() {
		r := &OffsetReader{
			Reader: strings.NewReader("0123456789"),
			Skip:   2,
			Stop:   8,
		}

		Convey("Read a byte", func() {
			p := make([]byte, 1)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 1)
			So(p, ShouldResemble, []byte("2"))
			So(r.Current(), ShouldEqual, 3)
		})

		Convey("Read all bytes", func() {
			p := make([]byte, 6)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 6)
			So(p, ShouldResemble, []byte("234567"))
			So(r.Current(), ShouldEqual, 8)

			Convey("Read 1 byte when all bytes already read", func() {
				p := make([]byte, 1)
				_, err := r.Read(p)
				So(err, ShouldEqual, io.EOF)
				So(r.Current(), ShouldEqual, 8)
			})
		})

		Convey("Read 7 bytes", func() {
			p := make([]byte, 7)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 6)
			So(p, ShouldResemble, append([]byte("234567"), 0))
			So(r.Current(), ShouldEqual, 8)
		})

	})

	Convey("Create a new offset reader with skip and stop base on not seekable reader", t, func() {
		r := &OffsetReader{
			Reader: bufio.NewReader(strings.NewReader("0123456789")),
			Skip:   2,
			Stop:   8,
		}

		Convey("Read a byte", func() {
			p := make([]byte, 1)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 1)
			So(p, ShouldResemble, []byte("2"))
			So(r.Current(), ShouldEqual, 3)
		})

		Convey("Read all bytes", func() {
			p := make([]byte, 6)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 6)
			So(p, ShouldResemble, []byte("234567"))
			So(r.Current(), ShouldEqual, 8)

			Convey("Read 1 byte when all bytes already read", func() {
				p := make([]byte, 1)
				_, err := r.Read(p)
				So(err, ShouldEqual, io.EOF)
				So(r.Current(), ShouldEqual, 8)
			})
		})

		Convey("Read 7 bytes", func() {
			p := make([]byte, 7)
			n, err := r.Read(p)
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 6)
			So(p, ShouldResemble, append([]byte("234567"), 0))
			So(r.Current(), ShouldEqual, 8)
		})

	})

	Convey("Create a new offset reader base on a faulty reader", t, func() {
		r := &OffsetReader{
			Reader: &FaultyReader{},
			Skip:   2,
		}

		Convey("Read 1 byte", func() {
			p := make([]byte, 1)
			_, err := r.Read(p)
			So(err, ShouldEqual, errRead)
		})
	})

	Convey("Create a new offset reader base on a faulty seek reader", t, func() {
		r := &OffsetReader{
			Reader: &FaultySeekableReader{},
			Skip:   2,
		}

		Convey("Read 1 byte", func() {
			p := make([]byte, 1)
			_, err := r.Read(p)
			So(err, ShouldEqual, errSeek)
		})
	})

}

var errSeek = errors.New("seek error")

type FaultySeekableReader struct {
	io.Reader
}

func (r *FaultySeekableReader) Read(p []byte) (n int, err error) {
	return r.Reader.Read(p)
}

func (r *FaultySeekableReader) Seek(offset int64, whence int) (int64, error) {
	return 0, errSeek
}

var errRead = errors.New("read error")

type FaultyReader struct {
}

func (r *FaultyReader) Read(p []byte) (n int, err error) {
	return 0, errRead
}
