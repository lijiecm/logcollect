package store

import (
	"errors"
	"github.com/weihualiu/logcollect/util"
)

//数据结构定义

type PackCommon struct {
	Header  byte
	Len     uint32 //4bytes
	Type    byte   //数据类型 0 heartbeat 1 api 2 log
	TagNum  uint8
	TagList []*Tag
	Date    []byte //7bytes
	Body    []byte
	Tail    byte
}

type Tag struct {
	Name []byte //12bytes
}

func NewPackCommon(data []byte) (*PackCommon, error) {
	packComm := new(PackCommon)
	packComm.Header = byte(0xF0)
	packComm.Tail = byte(0xFE)

	if data[0] != packComm.Header || data[len(data)-1] != packComm.Tail {
		return nil, errors.New("data struct parse failed, err dta header!")
	}

	packComm.Len = util.BytesToUInt32(data[1:5])
	if packComm.Len != uint32(len(data)) {
		return nil, errors.New("data struct parse failed, package length failed!")
	}
	packComm.Type = byte(data[5])
	//拆解标签
	packComm.TagNum = uint8(data[6])
	packComm.TagList = make([]*Tag, int(packComm.TagNum))
	for i := 0; i < int(packComm.TagNum); i++ {
		t := new(Tag)
		t.Name = make([]byte, 12)
		copy(t.Name, util.BytesTrim(data[7+12*i:7+12*(i+1)]))
		packComm.TagList[i] = t
	}

	beforeLen := 7 + int(packComm.TagNum)*12
	packComm.Date = data[beforeLen : beforeLen+8]
	packComm.Body = data[beforeLen+8 : len(data)-1]

	return packComm, nil

}
