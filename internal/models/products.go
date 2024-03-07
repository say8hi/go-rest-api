package models

type Product struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Price       float64  `json:"price"`
    Categories  []Category    `json:"categories"`
}


type CreateProductRequest struct {
    Name        string   `json:"name"`
    Description string   `json:"description"`
    Price       float64  `json:"price"`
    Categories  []int    `json:"categories"`
}

type ProductUpdateRequest struct {
    Name        *string   `json:"name,omitempty"`
    Description *string   `json:"description,omitempty"`
    Price       *float64  `json:"price,omitempty"`
    Categories  []int    `json:"categories,omitempty"`
}
