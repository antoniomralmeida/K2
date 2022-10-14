package classes

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type AttributeType byte

const (
	Trace AttributeType = iota
	Varchar
	Date
	Number
)

type Attribute struct {
	Name  string        `bson:"name"`
	Atype AttributeType `bson:"atype"`
	Value string        `bson:"value"`
}

type Class struct {
	Id         string      `bson:"_id"`
	Name       string      `bson:"name"`
	Icon       string      `bson:"icon"`
	Attributes []Attribute `bson:"attributes"`
}

type Command struct {
	text string
	bin  []*Token
}

type Instance struct {
	class *Class
	value []Attribute
}

type KB struct {
	ebnf       *EBNF
	kb_classes []*Class
}

func (e *KB) ReadBK(uri string, db string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	class := client.Database(db).Collection("Class")
	cur, err := class.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		var c Class
		cur.Decode(&c)
		fmt.Println(c)
	}
	/* var a []Attribute
	a = append(a, Attribute{"nome", Varchar, ""})
	x := Class{"Teste", "tese.jpg", a}
	result, err := class.InsertOne(ctx, x)
	if err != nil {
		log.Fatal(err)
	} else {
		newID := result.InsertedID
		fmt.Println(newID)
	}
	*/
	cur.Close(ctx)

}
