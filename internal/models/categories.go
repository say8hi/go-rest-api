
package models

type Category struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}


type CreateCategoryRequest struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}
