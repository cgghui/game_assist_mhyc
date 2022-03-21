package mhyc

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type ResIndex struct {
	Name   string
	Begin  int
	Length int
	Data   []byte
}

type Res struct {
	Index []ResIndex
	fp    *os.File
}

func (r *Res) GetData(i *ResIndex) error {
	_, _ = r.fp.Seek(0, io.SeekStart)
	_, _ = r.fp.Seek(int64(i.Begin), io.SeekCurrent)
	data := make([]byte, i.Length, i.Length)
	if _, err := r.fp.Read(data); err != nil {
		if err == io.EOF {
			err = nil
		}
		return err
	}
	i.Data = data
	return nil
}

func LoadDataRes(filePath string) (*Res, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	_, _ = f.Seek(0, io.SeekStart)
	_, _ = f.Seek(2, io.SeekCurrent)

	result := &Res{
		Index: make([]ResIndex, 0),
		fp:    f,
	}

	var s int16
	var v int32
	var l int
	data2 := make([]byte, 2, 2)
	for {
		if _, err = f.Read(data2); err != nil {
			if err == io.EOF {
				err = nil
			}
			return result, err
		}
		if bytes.Compare(data2, []byte{91, 10}) == 0 {
			return result, nil
		}
		if err = binary.Read(bytes.NewBuffer(data2), binary.BigEndian, &s); err != nil {
			return result, err
		}
		l = int(s)
		x := make([]byte, l+8, l+8)
		if _, err = f.Read(x); err != nil {
			if err == io.EOF {
				err = nil
			}
			return result, err
		}
		idx := ResIndex{Name: string(x[:l])}
		//
		v = 0
		if err = binary.Read(bytes.NewBuffer(x[l:4+l]), binary.BigEndian, &v); err != nil {
			return result, err
		}
		idx.Begin = int(v)
		//
		v = 0
		if err = binary.Read(bytes.NewBuffer(x[l+4:]), binary.BigEndian, &v); err != nil {
			return result, err
		}
		idx.Length = int(v)
		result.Index = append(result.Index, idx)
	}
}

type CfgClientMsg struct {
	Id     int    `json:"Id"`
	MSG    string `json:"MSG"`
	Module string `json:"Module"`
}
