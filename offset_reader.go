package xio

import "io"

var _ io.Reader = (*OffsetReader)(nil)

type OffsetReader struct {
	io.Reader

	// skip n bytes before read
	Skip int64

	// stop read after n bytes
	Stop int64

	// current offset
	cur int64
}

func (r *OffsetReader) Read(p []byte) (n int, err error) {
	if r.Stop != 0 && r.cur >= r.Stop {
		return 0, io.EOF
	}

	if r.Skip != 0 && r.cur < r.Skip {
		if seeker, ok := r.Reader.(io.Seeker); ok {
			// make sure seek and read behave the same
			_, err = seeker.Seek(r.Skip, io.SeekCurrent)
			if err != nil {
				return 0, err
			}
			r.cur = r.Skip
		} else {
			// read and discard
			for r.Skip-r.cur > int64(len(p)) {
				n, err = r.Reader.Read(p)
				r.cur += int64(n)
				if err != nil {
					return 0, err
				}
			}
			if r.Skip != r.cur {
				n, err = r.Reader.Read(p[:r.Skip-r.cur])
				r.cur += int64(n)
				if err != nil {
					return 0, err
				}
			}
		}
	}

	if r.Stop != 0 && r.cur+int64(len(p)) > r.Stop {
		p = p[:r.Stop-r.cur]
	}
	n, err = r.Reader.Read(p)
	r.cur += int64(n)
	return n, err
}

func (r *OffsetReader) Current() (n int64) {
	return r.cur
}
