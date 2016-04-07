package models

import "fmt"

type FileBody struct {
	Name    string    `json:"name,omitempty"`
	Part    int    `json:"part,omitempty"`
	Content []byte `json:"content,omitempty"`
}

func (self *FileBody) ToString() string {
	return fmt.Sprintf(
		"{Name: %v, Part: %v, Content: %v}",
		self.Name,
		self.Part,
		self.Content,
	)
}
