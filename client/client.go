package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "movie-ticket-booking/book"

	"google.golang.org/grpc"
)

const (
	address = "localhost: 5000"
)

type MyMovie struct {
	title       string
	genre       string
	year        int64
	showTime    string
	room        int64
	ticketsLeft int64
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	grpcClient := pb.NewBookerClient(conn)

	criterion := BrowseCriterion()
	switch criterion {
	case 1:
		title := GetMovie(grpcClient)
		order := BrowseMovies(grpcClient, title)
		HandleOrder(grpcClient, order)
	case 2:
		genre := GetGenre(grpcClient)
		order := BrowseGenre(grpcClient, genre)
		HandleOrder(grpcClient, order)
	case 3:
		showTime := GetShowTime(grpcClient)
		order := BrowseShowTime(grpcClient, showTime)
		HandleOrder(grpcClient, order)
	}

}

// BrowseOption asks for criterion to browse movies
func BrowseCriterion() int {
	fmt.Println("*-----* MOVIE TICKET BOOKING *-----*")
	fmt.Println()
	fmt.Println("How would you like to browse movies?")
	fmt.Println("1. By movie title\t2. By movie genre\t3. By showtime")
	fmt.Print("--> Your choice: ")
	var crit int
	fmt.Scanf("%d", &crit)
	return crit
}

// GetMovie ask for the movie user want to book
func GetMovie(grpcClient pb.BookerClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := grpcClient.ListMovie(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("GetMovie: %v", err)
	}
	// options slice holds data in the returned stream
	var options []string
	fmt.Println("Movies available at the theater:")
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("GetMovie: %v", err)
		}
		options = append(options, reply.GetOption())
	}
	return getOption(options)
}

// GetGenre ask for the genre user want to browse
func GetGenre(grpcClient pb.BookerClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := grpcClient.ListGenres(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("GetGenre: %v", err)
	}
	// options slice holds data in the returned stream
	var options []string
	fmt.Println("Movie genres at the theater:")
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("GetGenre: %v", err)
		}
		options = append(options, reply.GetOption())
	}
	return getOption(options)
}

// GetShowTime ask for the showtime user want to book
func GetShowTime(grpcClient pb.BookerClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := grpcClient.ListShowTimes(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("GetShowTime: %v", err)
	}
	// options slice holds data in the returned stream
	var options []string
	fmt.Println("Showtimes at the theater:")
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("GetShowTime: %v", err)
		}
		options = append(options, reply.GetOption())
	}
	return getOption(options)
}

// getOption shows all the browsing options and return one
func getOption(options []string) string {
	for i, option := range options {
		fmt.Printf("%d. %-15s", i+1, option)
	}
	fmt.Println()
	fmt.Print("--> Your choice: ")
	var choice int
	fmt.Scanf("%d", &choice)
	return options[choice-1]
}

func BrowseMovies(grpcClient pb.BookerClient, title string) MyMovie {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := grpcClient.MoviesByTitle(ctx, &pb.QueryParam{Param: title})
	if err != nil {
		log.Fatalf("BrowseMovie: %v", err)
	}
	// movies slice holds data in the returned stream
	var movies []MyMovie
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("BrowseMovie: %v", err)
		}
		var movie MyMovie
		movie.title = reply.Title
		movie.genre = reply.Genre
		movie.year = reply.Year
		movie.showTime = reply.ShowTime
		movie.room = reply.Room
		movie.ticketsLeft = reply.TicketsLeft
		movies = append(movies, movie)
	}
	return getBooking(movies)
}

func BrowseGenre(grpcClient pb.BookerClient, genre string) MyMovie {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := grpcClient.MoviesByGenre(ctx, &pb.QueryParam{Param: genre})
	if err != nil {
		log.Fatalf("BrowseGenre: %v", err)
	}
	// movies slice holds data in the returned stream
	var movies []MyMovie
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("BrowseGenre: %v", err)
		}
		var movie MyMovie
		movie.title = reply.Title
		movie.genre = reply.Genre
		movie.year = reply.Year
		movie.showTime = reply.ShowTime
		movie.room = reply.Room
		movie.ticketsLeft = reply.TicketsLeft
		movies = append(movies, movie)
	}
	return getBooking(movies)
}

func BrowseShowTime(grpcClient pb.BookerClient, showTime string) MyMovie {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := grpcClient.MoviesByTime(ctx, &pb.QueryParam{Param: showTime})
	if err != nil {
		log.Fatalf("BrowseShowTime: %v", err)
	}
	// movies slice holds data in the returned stream
	var movies []MyMovie
	for {
		reply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("BrowseShowTime: %v", err)
		}
		var movie MyMovie
		movie.title = reply.Title
		movie.genre = reply.Genre
		movie.year = reply.Year
		movie.showTime = reply.ShowTime
		movie.room = reply.Room
		movie.ticketsLeft = reply.TicketsLeft
		movies = append(movies, movie)
	}
	return getBooking(movies)
}

//getBooking show all the booking options and return one
func getBooking(movies []MyMovie) MyMovie {
	fmt.Print("\n\t\tHere we have some options for you:\n\n")
	fmt.Printf("%-8s|%-15s|%-15s|%-8s|%-15s|%-8s|%-8s\n", "Option", "Movie", "Genre", "Year", "Time", "Room", "Ticket left")
	fmt.Println("--------|---------------|---------------|--------|---------------|--------|------------")
	for i, movie := range movies {
		fmt.Printf("%-8d|%-15s|%-15s|%-8d|%-15v|%-8d|%-8d\n", i+1, movie.title, movie.genre, movie.year, movie.showTime, movie.room, movie.ticketsLeft)
	}
	fmt.Print("--> Your option: ")
	var option int
	fmt.Scanf("%d", &option)
	return movies[option-1]
}

// HandleOrder send booking request to server and handle the error returned
func HandleOrder(grpcClient pb.BookerClient, movie MyMovie) {
	fmt.Print("--> Number of tickets you would like to buy: ")
	var tickets int64
	fmt.Scanf("%d", &tickets)
	if tickets > movie.ticketsLeft {
		log.Fatalf("Sorry! Only %d ticket(s) left for your current option!", movie.ticketsLeft)
	}
	movieRequest := pb.Movie{
		Title:       movie.title,
		Genre:       movie.genre,
		Year:        movie.year,
		ShowTime:    movie.showTime,
		Room:        movie.room,
		TicketsLeft: movie.ticketsLeft,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := grpcClient.BuyTicket(ctx, &pb.BuyRequest{Movie: &movieRequest, Amount: tickets})
	if err != nil {
		log.Fatalf("HandleOrder: %v", err)
	}
	fmt.Println("Booking succeed! Congratulation!")
	fmt.Println("--- Here is your booking infomation:")
	fmt.Printf("--- Movie: %s (%d)\n", movie.title, movie.year)
	fmt.Printf("--- Genre: %s\n", movie.genre)
	fmt.Printf("--- Room: %d\t--- Time: %s\n", movie.room, movie.showTime)
	fmt.Println("--- Thank you! See you soon!")
}
