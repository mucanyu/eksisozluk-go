package model

// Topic fields
type Topic struct {
	Title         string
	NewEntryCount string
}

// Entry fields
type Entry struct {
	Text   string
	Author string
	Date   string
}

// GundemParams fields
type GundemParams struct {
	Kategori string
	Limit    int
	Page     int
}

// BaslikParams fields
type BaslikParams struct {
	Topic  string
	Page   int
	Limit  int
	Sukela bool
}
