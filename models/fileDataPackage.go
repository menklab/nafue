package models

import "fmt"

type FileDataPackage struct {
	Name    string    `json:"name,omitempty"`
	Content []byte `json:"content,omitempty"`
}

func (self *FileDataPackage) ToString() string {
	return fmt.Sprintf(
		"{Name: %v, Content: %v}",
		self.Name,
		self.Content,
	)
}
