package model

type Person struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" binding:"required" db:"name"`
	Surname     string `json:"surname" binding:"required" db:"surname"`
	Patronymic  string `json:"patronymic,omitempty" db:"patronymic"`
	Age         int    `json:"age,omitempty" db:"age"`
	Gender      string `json:"gender,omitempty" db:"gender"`
	Nationality string `json:"nationality,omitempty" db:"nationality"`
}
