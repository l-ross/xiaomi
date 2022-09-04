package miio

import (
	"encoding/binary"

	"github.com/ViRb3/slicewriteseek"
)

type data struct {
	sws *slicewriteseek.SliceWriteSeeker
}

func newData() *data {
	return &data{
		sws: slicewriteseek.New(),
	}
}

func newDataFromByte(b []byte) *data {
	return &data{
		sws: &slicewriteseek.SliceWriteSeeker{
			Buffer: b,
		},
	}
}

func (d *data) seek(pos uint16) {
	d.sws.Seek(int64(pos), 0)
}

func (d *data) bytes() []byte {
	return d.sws.Buffer
}

func (d *data) writeMagic() {
	d.write([]byte{0x21, 0x31})
}

func (d *data) writeUint16(v uint16) {
	binary.Write(d.sws, binary.BigEndian, v)
}

func (d *data) writeUint32(v uint32) {
	binary.Write(d.sws, binary.BigEndian, v)
}

func (d *data) write(b []byte) {
	d.sws.Write(b)
}

func (d *data) readInt16() (int16, error) {
	var v int16
	err := binary.Read(d.sws, binary.BigEndian, &v)
	return v, err
}

func (d *data) readBytes(l int32) ([]byte, error) {
	v := make([]byte, l)
	_, err := d.sws.Read(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}
