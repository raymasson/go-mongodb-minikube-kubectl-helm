package database

import (
	"crypto/tls"
	"fmt"

	"github.com/raymasson/go-mongodb-minikube-kubectl-helm/api/config"
	mgo "gopkg.in/mgo.v2"

	"net"
)

const (
	collectionName         = "person"
	databaseConnectionName = "admin"
	databaseCollectionName = "person"
)

// Manager interface
type Manager interface {
	ExecuteQuery(s func(Collection) error) error
}

// Session interface
type Session interface {
	DB(name string) DataLayer
	Close()
	Refresh()
}

// DataLayer interface
type DataLayer interface {
	C(name string) Collection
}

// Collection interface
type Collection interface {
	Find(query interface{}) Query
	Insert(docs ...interface{}) error
}

// Query interface
type Query interface {
	All(result interface{}) error
}

type manager struct {
	Session
}

// MongoSession is the Mongo session
type MongoSession struct {
	*mgo.Session
}

// MongoDatabase is the Mongo database
type MongoDatabase struct {
	*mgo.Database
}

// MongoCollection is the Mongo collection
type MongoCollection struct {
	*mgo.Collection
}

// MongoQuery is the Mongo query
type MongoQuery struct {
	*mgo.Query
}

// Find shadows *mgo.Collection to return a Query interface instead of *mgo.Query.
func (c MongoCollection) Find(query interface{}) Query {
	return &MongoQuery{Query: c.Collection.Find(query)}
}

// C shadows *mgo.DB to return a Collection interface instead of *mgo.Collection.
func (d MongoDatabase) C(name string) Collection {
	return &MongoCollection{Collection: d.Database.C(name)}
}

// DB shadows *mgo.DB to return a DataLayer interface instead of *mgo.Database.
func (s MongoSession) DB(name string) DataLayer {
	return &MongoDatabase{Database: s.Session.DB(name)}
}

// NewManager ...
func NewManager(s Session) Manager {
	return manager{s}
}

var (
	mgoSession *Session
	mongoURL   = "mongodb://%s:%s@%s/%s"
)

// Connect : opens a connection to the mongo DB
func Connect() Session {
	if mgoSession == nil {
		var connectionString = fmt.Sprintf(mongoURL, config.DbUser, config.DbPassword, config.DbHost, config.DbURI)

		dialInfo, err := mgo.ParseURL(connectionString)
		if err != nil {
			return MongoSession{}
		}

		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{InsecureSkipVerify: true})
		}

		mgoSession, err := mgo.DialWithInfo(dialInfo)
		if err != nil {
			return MongoSession{}
		}

		mgoSession.SetSafe(&mgo.Safe{})
		return MongoSession{mgoSession}
	}

	return *mgoSession
}

// Close : closes the mongo DB connection
func (s MongoSession) Close() {
	s.Session.Close()
}

// Refresh : refreshes the mongo DB connection
func (s MongoSession) Refresh() {
	s.Session.Refresh()
}

// ExecuteQuery : executes a query on the mongo DB
func (m manager) ExecuteQuery(s func(Collection) error) error {
	m.Session.Refresh()
	c := m.Session.DB(databaseCollectionName).C(collectionName)
	return s(c)
}
