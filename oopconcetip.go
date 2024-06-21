package main

import "fmt"

// Person struct with encapsulated fields
type Person struct {
	name string
	age  int
}

// NewPerson is a constructor for Person
func NewPerson(name string, age int) *Person {
	return &Person{name: name, age: age}
}

// Method to get the name
func (p *Person) GetName() string {
	return p.name
}

// Method to set the name
func (p *Person) SetName(name string) {
	p.name = name
}

// Method to get the age
func (p *Person) GetAge() int {
	return p.age
}

// Method to set the age
func (p *Person) SetAge(age int) {
	p.age = age
}

func main() {
	p := NewPerson("Alice", 30)
	fmt.Println(p.GetName()) // Outputs: Alice
	p.SetAge(31)
	fmt.Println(p.GetAge()) // Outputs: 31
}
