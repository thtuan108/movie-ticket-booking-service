package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "movie-ticket-booking/book"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	port = 5000
)

var db *gorm.DB

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
	dsn := os.Getenv("dbUser") + ":" + os.Getenv("dbPass") + "@tcp(127.0.0.1:3306)/theater?charset=utf8mb4&parseTime=true&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected!")
	}
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
func (s *server) ListMovie(ctx context.Context, in *pb.NoParam) (*pb.Options, error) {
	// query all the movie titles in db
	rows, err := db.Table("movies").Select("title", "year").Rows()
	if err != nil {
		return nil, fmt.Errorf("ListMovies: %v", err)
	}
	defer rows.Close()
	// titleList slice holds data returned from the query
	var titleList []string
	for rows.Next() {
		var title string
		var year int
		if err := rows.Scan(&title, &year); err != nil {
			return nil, fmt.Errorf("ListMovies: %v", err)
		}
		titleList = append(titleList, fmt.Sprintf("%s(%d)", title, year))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListMovies: %v", err)
	}
	return &pb.Options{List: titleList}, nil
}

// ListGenres returns options that correspond with movie genres
func (s *server) ListGenres(ctx context.Context, in *pb.NoParam) (*pb.Options, error) {
	// query all the movie genres in db
	rows, err := db.Table("movies").Select("genre").Rows()
	if err != nil {
		return nil, fmt.Errorf("ListGenes: %v", err)
	}
	defer rows.Close()
	// genreList slice holds data returned from the query
	var genreList []string
	for rows.Next() {
		var genre string
		if err := rows.Scan(&genre); err != nil {
			return nil, fmt.Errorf("ListGenres: %v", err)
		}
		genreList = append(genreList, genre)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListGenres: %v", err)
	}
	return &pb.Options{List: genreList}, nil
}

// List returns options that correspond with movie titles
func (s *server) ListShowTimes(ctx context.Context, in *pb.NoParam) (*pb.Options, error) {
	// query all the showtimes in db
	rows, err := db.Table("schedule").Distinct("show_time").Rows()
	if err != nil {
		return nil, fmt.Errorf("ListShowTimes: %v", err)
	}
	defer rows.Close()
	// timeList slice holds data returned from the query
	var timeList []string
	for rows.Next() {
		var showTime rawTime
		if err := rows.Scan(&showTime); err != nil {
			return nil, fmt.Errorf("ListShowTimes: %v", err)
		}
		showTimeStr, err := showTime.String()
		if err != nil {
			return nil, fmt.Errorf("ListShowTimes: %v", err)
		}
		timeList = append(timeList, showTimeStr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ListShowTimes: %v", err)
	}
	return &pb.Options{List: timeList}, nil
}

// MoviesByTitle queries for movies that have specified title
func (s *server) MoviesByTitle(ctx context.Context, in *pb.QueryParam) (*pb.Movies, error) {
	// get title to invoke query
	title := in.GetContent()
	// eliminate the year from the string (i.e: 'Soul(2020)' => 'Soul')
	title = title[:strings.IndexByte(title, '(')]
	rows, err := db.Table("movies").Select("movies.title", "movies.genre", "movies.year", "schedule.show_time", "schedule.room", "rooms.class", "rooms.seats-schedule.tickets_sold").Joins("JOIN schedule ON movies.id = schedule.movie_id").Joins("JOIN rooms ON rooms.id = schedule.room").Where("title = ?", title).Rows()
	if err != nil {
		return nil, fmt.Errorf("MoviesByTitle %s: %v", title, err)
	}
	defer rows.Close()
	// movieList slice holds data returned from rows
	var movieList []*pb.Movie
	for rows.Next() {
		var mov pb.Movie
		var r rawTime
		if err := rows.Scan(&mov.Title, &mov.Genre, &mov.Year, &r, &mov.Room, &mov.Class, &mov.TicketsLeft); err != nil {
			return nil, fmt.Errorf("MoviesByTitle %s: %v", title, err)
		}
		mov.ShowTime, err = r.String()
		if err != nil {
			return nil, fmt.Errorf("MoviesByTitle %s: %v", title, err)
		}
		movieList = append(movieList, &mov)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("MoviesByTitle %s: %v", title, err)
	}
	return &pb.Movies{Movie: movieList}, nil
}

// MoviesByGenre queries for movies that have specified genre
func (s *server) MoviesByGenre(ctx context.Context, in *pb.QueryParam) (*pb.Movies, error) {
	// get genre to invoke query
	genre := in.GetContent()
	rows, err := db.Table("movies").Select("movies.title", "movies.genre", "movies.year", "schedule.show_time", "schedule.room", "rooms.class", "rooms.seats-schedule.tickets_sold").Joins("JOIN schedule ON movies.id = schedule.movie_id").Joins("JOIN rooms ON rooms.id = schedule.room").Where("genre = ?", genre).Rows()
	if err != nil {
		return nil, fmt.Errorf("MoviesByGenre %s: %v", genre, err)
	}
	defer rows.Close()
	// movieList slice holds data returned from rows
	var movieList []*pb.Movie
	for rows.Next() {
		var mov pb.Movie
		var r rawTime
		if err := rows.Scan(&mov.Title, &mov.Genre, &mov.Year, &r, &mov.Room, &mov.Class, &mov.TicketsLeft); err != nil {
			return nil, fmt.Errorf("MoviesByGenre %s: %v", genre, err)
		}
		mov.ShowTime, err = r.String()
		if err != nil {
			return nil, fmt.Errorf("MoviesByGenre %s: %v", genre, err)
		}
		movieList = append(movieList, &mov)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("MoviesByGenre %s: %v", genre, err)
	}
	return &pb.Movies{Movie: movieList}, nil
}

// MoviesByTime queries for movies that start at specified time
func (s *server) MoviesByTime(ctx context.Context, in *pb.QueryParam) (*pb.Movies, error) {
	// get showTime to invoke query
	showTime := in.GetContent()
	rows, err := db.Table("movies").Select("movies.title", "movies.genre", "movies.year", "schedule.show_time", "schedule.room", "rooms.class", "rooms.seats-schedule.tickets_sold").Joins("JOIN schedule ON movies.id = schedule.movie_id").Joins("JOIN rooms ON rooms.id = schedule.room").Where("show_time = ?", showTime).Rows()
	if err != nil {
		return nil, fmt.Errorf("MoviesByTime %s: %v", showTime, err)
	}
	defer rows.Close()
	// movieList slice holds data returned from rows
	var movieList []*pb.Movie
	for rows.Next() {
		var mov pb.Movie
		var r rawTime
		if err := rows.Scan(&mov.Title, &mov.Genre, &mov.Year, &r, &mov.Room, &mov.Class, &mov.TicketsLeft); err != nil {
			return nil, fmt.Errorf("MoviesByTitle %s: %v", showTime, err)
		}
		mov.ShowTime, err = r.String()
		if err != nil {
			return nil, fmt.Errorf("MoviesByTitle %s: %v", showTime, err)
		}
		movieList = append(movieList, &mov)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("MoviesByTime %s: %v", showTime, err)
	}
	return &pb.Movies{Movie: movieList}, nil
}

// BuyTicket updates the number of tickets left in database base on amount of tickets bought
func (s *server) BuyTicket(ctx context.Context, buy *pb.BuyRequest) (*pb.NoParam, error) {
	movie := buy.GetMovie()
	amount := buy.GetAmount()
	if amount < 1 {
		return &pb.NoParam{}, fmt.Errorf("order failed: invalid amount of ticket(s)")
	}
	// return error if the amount of tickets to buy is greater than the amount of tickets left
	if amount > movie.TicketsLeft {
		return &pb.NoParam{}, fmt.Errorf("order failed: want %d ticket(s) but %d ticket(s) left", amount, movie.TicketsLeft)
	}
	// update the amount of tickets sold in database
	row := db.Table("movies").Select("id").Where("title = ?", movie.GetTitle()).Row()
	var movieId string
	err := row.Scan(&movieId)
	if err != nil {
		return &pb.NoParam{}, fmt.Errorf("BuyTicket: %v", err)
	}
	db.Table("schedule").Where("show_time = ? AND movie_id = ?", movie.GetShowTime(), movieId).Update("tickets_sold", movie.TicketsLeft-amount)
	// _, err := db.Exec("UPDATE schedule SET tickets_sold = tickets_sold+? WHERE show_time = ? AND movie_id = (SELECT id FROM movies WHERE title = ?)", amount, movie.GetShowTime(), movie.GetTitle())
	return &pb.NoParam{}, nil
}
