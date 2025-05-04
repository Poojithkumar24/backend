package mongo

/*schema:

users

{
  _id, name, email, password_hash, bookings: [booking_ids]
}

movies
{
  _id, title, genre, duration, description, screen, shows[]
}

shows
{
  _id, movie_id, date, time, screen_id, price_per_seat
}

screens
{
  _id, theater_name, seat_layout: [{ seat_no: "A1", type: "VIP" }]
}

bookings
{
  _id, user_id, show_id, seats: ["A1", "A2"], status: "confirmed", total_price
}
*/
