package models

import "fmt"

type FileDataPackage struct {
	Name    string    `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
}

func (self *FileDataPackage) ToString() string {
	return fmt.Sprintf(
		"{Name: %v, Content: %v}",
		self.Name,
		self.Content,
	)
}
