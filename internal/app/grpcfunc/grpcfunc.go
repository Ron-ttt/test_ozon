package grpcfunc

import (
	"context"
	"net/url"
	"testozon/internal/app/handlers"
	pb "testozon/internal/app/proto"
	"testozon/internal/app/storage"
	"testozon/internal/app/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var store storage.Storage

const localhost = "localhost:3200"
const baseURL = "http://" + localhost

type ShortenerServer struct {
	// нужно встраивать тип pb.Unimplemented<TypeName>
	// для совместимости с будущими версиями
	pb.UnimplementedShortenerServer
}

func Init() {
	dbAdress := "postgresql://postgres_user:postgres_password@postgres_container:5432/postgres?sslmode=disable"
	if handlers.Flag {
		dBStorage, err := storage.NewDBStorage(dbAdress)
		if err == nil {
			store = dBStorage
			return
		}
	}

	store = storage.NewMapStorage()
}

func (s *ShortenerServer) RedirectTo(ctx context.Context, in *pb.RedirectToRequest) (*pb.RedirectToResponse, error) {
	var response pb.RedirectToResponse

	originalURL, err := store.Get(in.ShortURL)
	if err != nil {
		return nil, status.Error(codes.Internal, "Err occurred getting url")
	}

	response.OriginalURL = originalURL
	return &response, nil
}

func (s *ShortenerServer) IndexPage(ctx context.Context, in *pb.IndexPageRequest) (*pb.IndexPageResponse, error) {
	var response pb.IndexPageResponse

	_, err1 := url.ParseRequestURI(in.OriginalUrl)
	if err1 != nil {
		err := status.Error(codes.InvalidArgument, "Invalid URL")
		return nil, err
	}

	oldShortURL, err := store.Find(in.OriginalUrl)
	if err == nil {
		response.ShortUrl = baseURL + "/" + oldShortURL
		return &response, nil
	}

	urlLen := 10
	code := utils.RandString(urlLen)
	err = store.Add(code, in.OriginalUrl)
	if err != nil {
		err = status.Error(codes.Internal, "Err occurred adding url")
		return nil, err
	}

	response.ShortUrl = code
	return &response, nil
}
