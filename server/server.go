package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "movie-ticket-booking/book"

	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

const (
	port = 5000
)

var db *sql.DB

type server struct {
	pb.UnimplementedBookerServer
}

type rawTime []byte

func (r rawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(r))
}

func (r rawTime) String() (string, error) {
	timeStr, err := r.Time()
	if err != nil {
		return "", fmt.Errorf("time format: %v", err)
	}
	return timeStr.Format("15:04:05"), nil
}

func main() {
	// Connect to database
	// Capture connection properties
	cfg := mysql.Config{
		User:      os.Getenv("dbUser"),
		Passwd:    os.Getenv("dbPass"),
		Net:       "tcp",
		Addr:      "127.0.0.1:3306",
		DBName:    "theater",
		ParseTime: true,
	}
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// start server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost: %d", port))
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBookerServer(grpcServer, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("fail to serve: %v", err)
	}
}

// ListMovies returns options that correspond with movie titles
func (s *server) ListMovie(in *pb.NoParam, stream pb.Booker_ListMovieServer) error {
	// query all the movie titles in db
	rows, err := db.Query("SELECT title FROM movies")
	if err != nil {
		return fmt.Errorf("ListMovies: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var title string
		if err := rows.Scan(&title); err != nil {
			return fmt.Errorf("ListMovies: %v", err)
		}
		if err := stream.Send(&pb.Options{Option: title}); err != nil {
			return fmt.Errorf("ListMovies: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("ListMovies: %v", err)
	}
	return nil
}

// ListGenres returns options that correspond with movie genres
func (s *server) ListGenres(in *pb.NoParam, stream pb.Booker_ListGenresServer) error {
	// query all the movie genres in db
	rows, err := db.Query("SELECT DISTINCT genre FROM movies")
	if err != nil {
		return fmt.Errorf("ListGenes: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return fmt.Errorf("ListGenres: %v", err)
		}
		if err := stream.Send(&pb.Options{Option: genre}); err != nil {
			return fmt.Errorf("ListMovies: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("ListGenres: %v", err)
	}
	return nil
}

// List returns options that correspond with movie titles
func (s *server) ListShowTimes(in *pb.NoParam, stream pb.Booker_ListShowTimesServer) error {
	// query all the showtimes in db
	rows, err := db.Query("SELECT DISTINCT show_time FROM schedule")
	if err != nil {
		return fmt.Errorf("ListShowTimes: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var showTime rawTime
		if err := rows.Scan(&showTime); err != nil {
			return fmt.Errorf("ListShowTimes: %v", err)
		}
		showTimeStr, err := showTime.String()
		if err != nil {
			return fmt.Errorf("ListShowTimes: %v", err)
		}
		if err := stream.Send(&pb.Options{Option: showTimeStr}); err != nil {
			return fmt.Errorf("ListShowTimes: %v", err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("ListShowTimes: %v", err)
	}
	return nil
}

// MoviesByTitle queries for movies that have specified title
func (s *server) MoviesByTitle(in *pb.QueryParam, stream pb.Booker_MoviesByTitleServer) error {
	// get title to invoke query
	title := in.GetParam()
	rows, err := db.Query("SELECT title, genre, year, show_time, room, (tickets-tickets_sold) FROM movies m JOIN schedule s ON m.id = s.movie_id WHERE title = ?", title)
	if err != nil {
		return fmt.Errorf("MoviesByTitle %s: %v", title, err)
	}
	defer rows.Close()
	// stream rows returned to client
	for rows.Next() {
		mov := &pb.Movie{}
		var r rawTime
		if err := rows.Scan(&mov.Title, &mov.Genre, &mov.Year, &r, &mov.Room, &mov.TicketsLeft); err != nil {
			return fmt.Errorf("MoviesByTitle %s: %v", title, err)
		}
		mov.ShowTime, err = r.String()
		if err != nil {
			return fmt.Errorf("MoviesByTitle %s: %v", title, err)
		}
		if err := stream.Send(mov); err != nil {
			return fmt.Errorf("MoviesByTitle %s: %v", title, err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("MoviesByTitle %s: %v", title, err)
	}
	return nil
}

// MoviesByGenre queries for movies that have specified genre
func (s *server) MoviesByGenre(in *pb.QueryParam, stream pb.Booker_MoviesByGenreServer) error {
	// get genre to invoke query
	genre := in.GetParam()
	rows, err := db.Query("SELECT title, genre, year, show_time, room, (tickets-tickets_sold) FROM movies m JOIN schedule s ON m.id = s.movie_id WHERE genre = ?", genre)
	if err != nil {
		return fmt.Errorf("MoviesByGenre %s: %v", genre, err)
	}
	defer rows.Close()
	// stream rows returned to client
	for rows.Next() {
		var mov pb.Movie
		var r rawTime
		if err := rows.Scan(&mov.Title, &mov.Genre, &mov.Year, &r, &mov.Room, &mov.TicketsLeft); err != nil {
			return fmt.Errorf("MoviesByGenre %s: %v", genre, err)
		}
		mov.ShowTime, err = r.String()
		if err != nil {
			return fmt.Errorf("MoviesByGenre %s: %v", genre, err)
		}
		if err := stream.Send(&mov); err != nil {
			return fmt.Errorf("MoviesByGenre %s: %v", genre, err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("MoviesByGenre %s: %v", genre, err)
	}
	return nil
}

// MoviesByTime queries for movies that start at specified time
func (s *server) MoviesByTime(in *pb.QueryParam, stream pb.Booker_MoviesByTimeServer) error {
	// get showTime to invoke query
	showTime := in.GetParam()
	rows, err := db.Query("SELECT title, genre, year, show_time, room, (tickets-tickets_sold) FROM movies m JOIN schedule s ON m.id = s.movie_id WHERE show_time = ?", showTime)
	if err != nil {
		return fmt.Errorf("MoviesByTime %s: %v", showTime, err)
	}
	defer rows.Close()
	// stream rows returned to client
	for rows.Next() {
		var mov pb.Movie
		var r rawTime
		if err := rows.Scan(&mov.Title, &mov.Genre, &mov.Year, &r, &mov.Room, &mov.TicketsLeft); err != nil {
			return fmt.Errorf("MoviesByTitle %s: %v", showTime, err)
		}
		mov.ShowTime, err = r.String()
		if err != nil {
			return fmt.Errorf("MoviesByTitle %s: %v", showTime, err)
		}
		if err := stream.Send(&mov); err != nil {
			return fmt.Errorf("MoviesByTime %s: %v", showTime, err)
		}
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("MoviesByTime %s: %v", showTime, err)
	}
	return nil
}

// BuyTicket updates the number of tickets left in database base on amount of tickets bought
func (s *server) BuyTicket(ctx context.Context, buy *pb.BuyRequest) (*pb.NoParam, error) {
	movie := buy.GetMovie()
	amount := buy.GetAmount()
	// return error if the amount of tickets to buy is greater than the amount of tickets left
	if amount > movie.TicketsLeft {
		return &pb.NoParam{}, fmt.Errorf("order failed: want %d ticket(s) but %d ticket(s) left", amount, movie.TicketsLeft)
	}
	// update the amount of tickets sold in database
	_, err := db.Exec("UPDATE schedule SET tickets_sold = tickets_sold+? WHERE show_time = ? AND movie_id = (SELECT id FROM movies WHERE title = ?)", amount, movie.GetShowTime(), movie.GetTitle())
	if err != nil {
		return &pb.NoParam{}, fmt.Errorf("BuyTicket: %v", err)
	}
	return &pb.NoParam{}, nil
}
