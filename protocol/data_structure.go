package protocol

import (
    "errors"
    "github.com/jmtp/jmtp-client-go/util"
    "bytes"
)

const (
    PacketMaxSize = 268435455
    PacketMinSize = 3
)

type Operation interface {
    GetBytes() []byte
}

type BaseDataType struct {
    message []byte
}

// Tiny String
type TinyString struct {
    BaseDataType
    original    string
}

func (ts *TinyString) SetContent(str string) error {

    content := []byte(str)
    contentLen := len(content)
    if content == nil || contentLen > 255 {
        return errors.New("invalidate parameter")
    }
    message, err := getMessage(content)
    if err != nil {
        return err
    }
    ts.message = message
    ts.original = str

    return nil
}

func (ts *TinyString) GetBytes() []byte {
    return ts.message
}

func (ts *TinyString) GetContent() string {
    return ts.original
}

// Tine bytes
type TinyBytes struct {
    BaseDataType
    original    []byte
}

func (tb *TinyBytes) SetContent(input []byte) error {

    if input == nil || len(input) > 255 || len(input) == 0{
        return errors.New("invalidate parameter")
    }
    message, err := getMessage(input)
    if err != nil {
        return err
    }
    tb.message = message
    tb.original = input

    return nil
}

func (tb *TinyBytes) GetBytes() []byte {
    return tb.message
}

func (tb *TinyBytes) GetContent() []byte {
    return tb.original
}

// short varchar
type ShortVarchar struct {
    BaseDataType
    original    string
}

func (sv *ShortVarchar) SetContent(input string) error {

    len := len(input)
    if len > 18385 || len <= 0 {
        return errors.New("remaining length overflow")
    }
    message, err := getMessage([]byte(input))
    if err != nil {
        return nil
    }
    sv.message = message
    sv.original = input

    return nil
}

func (sv *ShortVarchar) GetBytes() []byte {
    return sv.message
}

func (sv *ShortVarchar) GetContent() []byte {
    return sv.message
}

// ShortBytes
type ShortBytes struct {
    BaseDataType
    original    []byte
}

func (sv *ShortBytes) SetContent(input []byte) error {

    len := len(input)
    if len > 18385 || len <= 0 {
        return errors.New("remaining length overflow")
    }
    message, err := getMessage(input)
    if err != nil {
        return nil
    }
    sv.message = message
    sv.original = input

    return nil
}

func (sv *ShortBytes) GetBytes() []byte {
    return sv.message
}

func (sv *ShortBytes) GetContent() []byte {
    return sv.original
}

type Tuple struct {
    BaseDataType
    key *ShortVarchar
    val *ShortBytes
}

func (t *Tuple) SetVal(key string, val []byte) error {

    t.key = &ShortVarchar{}
    t.key.SetContent(key)
    t.val = &ShortBytes{}
    t.val.SetContent(val)

    len := len(t.key.GetBytes()) + len(t.val.GetBytes())
    t.message = make([]byte, len)
    t.message = append(t.message, t.key.GetBytes()...)
    t.message = append(t.message, t.val.GetBytes()...)

    return nil
}

func (t *Tuple) GetKey() string {
    return t.key.original
}

func (t *Tuple) GetVal() []byte {
    return t.val.original
}

func (t *Tuple) GetBytes() []byte {
    return t.message
}

// TinyMap
type TinyMap struct {
    BaseDataType
    tuples []*Tuple
}

func (ts *TinyMap) Set(key string, val []byte) error {
    for _, v := range ts.tuples {
        if v.key.original == key {
            return errors.New("contains same key in tiny map")
        }
    }
    tuple := &Tuple{}
    tuple.SetVal(key, val)
    ts.tuples = append(ts.tuples, tuple)

    return nil
}

func (ts *TinyMap) Get(key string) []byte {
    var out []byte
    for _, v := range ts.tuples {
        if v.key.original == key {
            out = v.val.original
            break
        }
    }
    return out
}

func (ts *TinyMap) GetBytes() ([]byte, error) {

    var buffer bytes.Buffer
    err := EncodeRemainingLength(uint(len(ts.tuples)), buffer)
    if err != nil {
        return nil, err
    }
    for _, t := range ts.tuples {
        buffer.Write(t.GetBytes())
    }

    return buffer.Bytes(), nil
}

// format message
func getMessage(content []byte) ([]byte, error) {

    var buffer bytes.Buffer
    contentLen := len(content)
    err := EncodeRemainingLength(uint(contentLen), buffer)
    if err != nil {
        return nil, err
    }
    buffer.Write(content)
    out := make([]byte, buffer.Len())
    num, err := buffer.Read(out)
    if err != nil {
        return out, err
    }

    return out[:num], nil
}

func EncodeRemainingLength(len uint, out bytes.Buffer) error {

    if len > PacketMaxSize {
        return errors.New("remaining length overflow")
    }
    x := len
    var encodeByte byte
    for {
        encodeByte = util.Uint2Byte(uint64(util.FloorMod(int(x), 128)))
        if x > 0 {
            out.WriteByte(encodeByte | 128)
        } else {
            out.WriteByte(encodeByte)
        }
        if x <= 0 {
            break
        }
    }

    return nil
}

func DecodeRemainingLength(in bytes.Buffer) (int, error) {

    multiplier := 1
    remainingLength := 0
    for {
        encodeByte, err := in.ReadByte()
        if err != nil {
            return remainingLength, err
        }
        i, err := util.Byte2Int(encodeByte & 127)
        if err != nil {
            return remainingLength, err
        }
        remainingLength += i * multiplier
        if (encodeByte & 128) != 0 {
            if multiplier == 128 * 128 * 128 {
                return remainingLength, errors.New("malformed remaining length")
            }
        }
    }

    return remainingLength, nil
}
