package models

type PetResponse struct {
    ID         int64      `json:"id"`
    Category   PetCategory   `json:"category"`
    Name       string     `json:"name"`
    PhotoUrls  []string   `json:"photoUrls"`
    Tags       []Tag      `json:"tags"`
    Status     string     `json:"status"`
}

type PetCategory struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

type Tag struct {
    ID   int64  `json:"id"`
    Name string `json:"name"`
}

type Product struct {
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    Price       float64  `json:"price,omitempty"`       
    Categories  []Category `json:"categories"`
}

type Category struct {
    Name string `json:"name"`
    Description string `json:"description"`
}
