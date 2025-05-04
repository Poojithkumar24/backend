package mongo

type User struct {
	UserID       string   `bson:"userID" json:"userID"`
	IsAdmin      bool     `bson:"isAdmin" json:"isAdmin"`
	Name         string   `bson:"name" json:"name"`
	Email        string   `bson:"email" json:"email"`
	PasswordHash string   `bson:"passwordHash" json:"passwordHash"`
	Bookings     []string `bson:"bookings" json:"bookings"`
}
type Movie struct {
	MovieID     string   `bson:"movieID,omitempty" json:"movieID"`
	Title       string   `bson:"title" json:"title"`
	Genre       string   `bson:"genre" json:"genre"`
	Duration    int      `bson:"duration" json:"duration"` // duration in minutes
	Description string   `bson:"description" json:"description"`
	Screen      string   `bson:"screen" json:"screen"` // screen name or ID
	Shows       []string `bson:"shows" json:"shows"`
}
type Show struct {
	SHowID       string  `bson:"showID,omitempty" json:"showID"`
	MovieID      string  `bson:"movie_id" json:"movieId"`
	Date         string  `bson:"date" json:"date"` // format: YYYY-MM-DD
	Time         string  `bson:"time" json:"time"` // format: HH:MM
	ScreenID     string  `bson:"screen_id" json:"screenId"`
	PricePerSeat float64 `bson:"pricePerSeat" json:"pricePerSeat"`
}
type Seat struct {
	SeatNo string `bson:"seatNo" json:"seatNo"`
	Type   string `bson:"type" json:"type"` // e.g., VIP, Standard
}

type Screen struct {
	ScreenID    string `bson:"screenID,omitempty" json:"screenID"`
	TheaterName string `bson:"theater_name" json:"theaterName"`
	SeatLayout  []Seat `bson:"seatLayout" json:"seatLayout"`
}
type Booking struct {
	BookingID  string   `bson:"bookingID,omitempty" json:"bookingID"`
	UserID     string   `bson:"userID" json:"userID"`
	ShowID     string   `bson:"showID" json:"showID"`
	Seats      []string `bson:"seats" json:"seats"`   // e.g., ["A1", "A2"]
	Status     string   `bson:"status" json:"status"` // e.g., "confirmed"
	TotalPrice float64  `bson:"totalPrice" json:"totalPrice"`
}
