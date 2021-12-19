package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "movie-ticket-booking/book"

	"google.golang.org/grpc"
)

const (
	address = "localhost: 5000"
)

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
	default:
		log.Fatal("Invalid choice!")
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

	options, err := grpcClient.ListMovie(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("GetMovie: %v", err)
	}
	return getOption(options.List)
}

// GetGenre ask for the genre user want to browse
func GetGenre(grpcClient pb.BookerClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options, err := grpcClient.ListGenres(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("GetGenre: %v", err)
	}
	return getOption(options.List)
}

// GetShowTime ask for the showtime user want to book
func GetShowTime(grpcClient pb.BookerClient) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options, err := grpcClient.ListShowTimes(ctx, &pb.NoParam{})
	if err != nil {
		log.Fatalf("GetShowTime: %v", err)
	}
	return getOption(options.List)
}

// getOption shows all the browsing options and return one
func getOption(options []string) string {
	for i, option := range options {
		fmt.Printf("%d. %s   ", i+1, option)
	}
	fmt.Println()
	fmt.Print("--> Your choice: ")
	var choice int
	fmt.Scanf("%d", &choice)
	if choice < 1 || choice > len(options) {
		log.Fatal("Invalid choices!")
	}
	return options[choice-1]
}

func BrowseMovies(grpcClient pb.BookerClient, title string) *pb.Movie {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movieList, err := grpcClient.MoviesByTitle(ctx, &pb.QueryParam{Content: title})
	if err != nil {
		log.Fatalf("BrowseMovie: %v", err)
	}
	return getBooking(movieList.Movie)
}

func BrowseGenre(grpcClient pb.BookerClient, genre string) *pb.Movie {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movieList, err := grpcClient.MoviesByGenre(ctx, &pb.QueryParam{Content: genre})
	if err != nil {
		log.Fatalf("BrowseGenre: %v", err)
	}
	return getBooking(movieList.Movie)
}

func BrowseShowTime(grpcClient pb.BookerClient, showTime string) *pb.Movie {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	movieList, err := grpcClient.MoviesByTime(ctx, &pb.QueryParam{Content: showTime})
	if err != nil {
		log.Fatalf("BrowseShowTime: %v", err)
	}
	return getBooking(movieList.Movie)
}

//getBooking show all the booking options and return one
func getBooking(movies []*pb.Movie) *pb.Movie {
	fmt.Print("\n\t\t\tHere we have some options for you:\n\n")
	fmt.Printf("%-8s|%-15s|%-15s|%-8s|%-15s|%-8s|%-12s|%-8s\n", "Option", "Movie", "Genre", "Year", "Time", "Room", "Class", "Tickets left")
	fmt.Println("--------|---------------|---------------|--------|---------------|--------|------------|------------")
	for i, movie := range movies {
		fmt.Printf("%-8d|%-15s|%-15s|%-8d|%-15v|%-8d|%-12s|%-8d\n", i+1, movie.Title, movie.Genre, movie.Year, movie.ShowTime, movie.Room, movie.Class, movie.TicketsLeft)
	}
	fmt.Print("--> Your option: ")
	var option int
	fmt.Scanf("%d", &option)
	if option < 1 || option > len(movies) {
		log.Fatal("Invalid option!")
	}
	return movies[option-1]
}

// HandleOrder send booking request to server and handle the error returned
func HandleOrder(grpcClient pb.BookerClient, movie *pb.Movie) {
	fmt.Print("--> Number of tickets you would like to buy: ")
	var tickets int64
	fmt.Scanf("%d", &tickets)
	if tickets < 1 {
		log.Fatal("Invalid input!")
	}
	if tickets > movie.TicketsLeft {
		log.Fatalf("Sorry! Only %d ticket(s) left for your current option!", movie.TicketsLeft)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := grpcClient.BuyTicket(ctx, &pb.BuyRequest{Movie: movie, Amount: tickets})
	if err != nil {
		log.Fatalf("HandleOrder: %v", err)
	}
	fmt.Println("Booking is successful! Congratulation!")
	fmt.Println("--- Here is your booking infomation:")
	fmt.Printf("--- Movie: %s (%d)\n", movie.Title, movie.Year)
	fmt.Printf("--- Genre: %s\n", movie.Genre)
	fmt.Printf("--- Room: %d\t--- Time: %s\n", movie.Room, movie.ShowTime)
	fmt.Println("--- Thank you! See you soon!")
}
