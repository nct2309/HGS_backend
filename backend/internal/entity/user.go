package entity

// var (
// 	ERR_USER_PASSWORD_NOT_MATCH = errors.New("Password not match")
// 	ERR_USER_NOT_FOUND          = errors.New("User not found")
// 	ERR_USER_FINE_EXCEED        = errors.New("User fine exceeds the limit")
// )

// // User from last project
// type User struct {
// 	ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
// 	Phonenum      string             `json:"phonenum" bson:"phonenum"`
// 	Age           int                `json:"age" bson:"age"`
// 	Gender        string             `json:"gender" bson:"gender"`
// 	SSN           string             `json:"ssn" bson:"ssn"`
// 	Name          string             `json:"name" bson:"name"`
// 	Role          int                `json:"flag" bson:"flag"` // user for identify 1 for member, 2 for librarian or manager is 3?
// 	CountFine     int                `json:"countfine" bson:"countfine"`
// 	Username      string             `json:"username" bson:"username"`
// 	Password      string             `json:"password" bson:"password"`
// 	ReservingList []UserActivity     `json:"reservinglist" bson:"reservinglist"`
// 	BorrowingList []UserActivity     `json:"borrowinglist" bson:"borrowinglist"`
// 	BorrowedList  []UserActivity     `json:"borrowedlist" bson:"borrowedlist"`
// }

// type UserActivity struct {
// 	BookId       primitive.ObjectID `json:"bookId" bson:"bookId"`
// 	BookName     string             `json:"bookName" bson:"bookName"`
// 	StartDate    time.Time          `json:"startDate" bson:"startDate"`
// 	EndDate      time.Time          `json:"endDate" bson:"endDate"`
// 	ExtendedDate time.Time          `json:"extendedDate" bson:"extendedDate"`
// 	Location     string             `json:"location" bson:"location"`
// }

// // type Member struct {
// // 	CountFine int `json:"countfine" bson:"countfine"`
// // }

// // type Librarian struct {
// // 	????
// // }

// // type Manager struct {
// // 	????
// // }

type User struct {
	ID       int    `gorm:"primaryKey;column:User_id" json:"user_id"`
	Username string `gorm:"username" json:"username"`
	Password string `gorm:"password" json:"password"`
}
