package persistency

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"
)

var lock sync.Mutex

var marshal = func(v proto.Message) (io.Reader, error) {
	imessagebytes, err := proto.Marshal(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(imessagebytes), nil
}

var unmarshal = func(r io.Reader, v proto.Message) error {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)
	err := proto.Unmarshal(buf.Bytes(), v)
	return err
}

//Save a proto message
func Save(path string, v proto.Message) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

//Load a Proto Message from path
func Load(path string, v proto.Message) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return unmarshal(f, v)
}
