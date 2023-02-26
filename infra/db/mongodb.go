package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/big"
	"reflect"
	"time"
)

type MongoDB struct {
	connString string
	client     *mongo.Client
}

var (
	tUUID       = reflect.TypeOf(uuid.UUID{})
	uuidSubtype = bsontype.BinaryUUID

	tBigInt = reflect.TypeOf(&big.Int{})
)

func NewMongoDb() *MongoDB {
	return new(MongoDB)
}

func (mongoDb *MongoDB) Connect(connString string) error {
	if mongoDb.IsConnected() {
		return fmt.Errorf("DB is already connected")
	}

	mongoRegistry := bson.NewRegistryBuilder().
		RegisterTypeEncoder(tUUID, bsoncodec.ValueEncoderFunc(uuidEncodeValue)).
		RegisterTypeDecoder(tUUID, bsoncodec.ValueDecoderFunc(uuidDecodeValue)).
		RegisterTypeEncoder(tBigInt, bsoncodec.ValueEncoderFunc(bigIntEncodeValue)).
		RegisterTypeDecoder(tBigInt, bsoncodec.ValueDecoderFunc(bigIntDecodeValue)).
		Build()

	opt := options.Client().ApplyURI(connString).SetRegistry(mongoRegistry)

	client, err := mongo.NewClient(opt)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoDb.client = client
	mongoDb.connString = connString

	return nil
}

func (mongoDb *MongoDB) IsConnected() bool {
	return mongoDb.client != nil
}

func (mongoDb *MongoDB) GetClient() interface{} {
	return mongoDb.client
}

func bigIntEncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tBigInt {
		return bsoncodec.ValueEncoderError{Name: "bigIntEncodeValue", Types: []reflect.Type{tBigInt}, Received: val}
	}
	b := val.Interface().(*big.Int)
	// stores big int as base16 (hex) strings
	return vw.WriteString(b.Text(16))
}

func bigIntDecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tBigInt {
		return bsoncodec.ValueDecoderError{Name: "bigIntDecodeValue", Types: []reflect.Type{tBigInt}, Received: val}
	}

	var data string
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.String:
		data, err = vr.ReadString()
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	bigInt, ok := new(big.Int).SetString(data, 16)
	if !ok {
		return fmt.Errorf("failed to convert base16 string %s to big int", data)
	}
	val.Set(reflect.ValueOf(bigInt))
	return nil
}

func uuidEncodeValue(_ bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	if !val.IsValid() || val.Type() != tUUID {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}
	b := val.Interface().(uuid.UUID)
	return vw.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(_ bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() || val.Type() != tUUID {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{tUUID}, Received: val}
	}

	var data []byte
	var subtype byte
	var err error
	switch vrType := vr.Type(); vrType {
	case bsontype.Binary:
		data, subtype, err = vr.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = vr.ReadNull()
	case bsontype.Undefined:
		err = vr.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}

	if err != nil {
		return err
	}
	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}
	val.Set(reflect.ValueOf(uuid2))
	return nil
}
