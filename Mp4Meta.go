package utils

import (
	"bytes"
	"encoding/binary"
	"io"
)

// BoxHeader 信息头
type BoxHeader struct {
	Size       uint32
	FourccType [4]byte
	Size64     uint64
}

// GetMP4Duration 获取视频时长，以秒计
func GetMP4Duration(reader io.ReaderAt) []byte {
	var err error
	var info = make([]byte, 0x10)
	var boxHeader BoxHeader
	var offset int64 = 0
	var saveData = make(map[string][]byte)
	// 获取moov结构偏移
	for {
		_, err = reader.ReadAt(info, offset)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		boxHeader = getHeaderBoxInfo(info)
		fourccType := getFourccType(boxHeader)
		_bytes := make([]byte, boxHeader.Size)
		if fourccType == "mdat" && boxHeader.Size == 1 {
			_bytes = make([]byte, boxHeader.Size64)
		}
		_, _ = reader.ReadAt(_bytes, offset)
		saveData[fourccType] = _bytes[:]
		if fourccType == "mdat" && boxHeader.Size == 1 {
			offset += int64(boxHeader.Size64)
			continue
		}
		offset += int64(boxHeader.Size)
	}
	var mD []byte

	mD = append(mD, saveData["ftyp"]...)
	mD = append(mD, saveData["free"]...)
	mD = append(mD, saveData["mdat"]...)
	mD = append(mD, saveData["moov"]...)
	//FilePutContentsbytes("test.ts",saveData["mdat"][8:])
	return mD
}

// getHeaderBoxInfo 获取头信息
func getHeaderBoxInfo(data []byte) (boxHeader BoxHeader) {
	buf := bytes.NewBuffer(data)
	_ = binary.Read(buf, binary.BigEndian, &boxHeader)
	return
}

// getFourccType 获取信息头类型
func getFourccType(boxHeader BoxHeader) (fourccType string) {
	fourccType = string(boxHeader.FourccType[:])
	return
}
