package model

import "sync"

type Server struct {
	Rooms map[string]*Room
	Lock sync.RWMutex
}
type Room struct {
	UserIdA,UserIdB int64

}
