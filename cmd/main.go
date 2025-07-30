package main

import (
	"fmt"

	gofiledb "github.com/Kaushik1766/GoFileDB"
)

type User struct {
	Id   string
	Name string
	Age  int
}

func (u User) GetID() string {
	return u.Id
}

func main() {
	var db gofiledb.Repository[User] = gofiledb.NewFileDBRepository[User]()
	db.Save(User{"adsfaddssdf", "kaushik", 41})
	// fmt.Println(db.GetByParameter(map[string]any{"Age": 41}))
	// db.DeleteByParameter(map[string]any{"Age": 41})
	fmt.Println(db.GetByParameter(map[string]any{}))
}
