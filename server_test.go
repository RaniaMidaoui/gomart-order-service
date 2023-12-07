package test

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"

	"github.com/RaniaMidaoui/gomart-order-service/pkg/client"
	"github.com/RaniaMidaoui/gomart-order-service/pkg/db"
	"github.com/RaniaMidaoui/gomart-order-service/pkg/pb"
	productServices "github.com/RaniaMidaoui/gomart-order-service/pkg/product_services"
	services "github.com/RaniaMidaoui/gomart-order-service/pkg/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func productServer(ctx context.Context) (pb.ProductServiceClient, func(), error) {

	lis := bufconn.Listen(1024 * 1024)

	s := grpc.NewServer()

	h := db.Mock()

	ss := productServices.Server{
		H: h,
	}

	pb.RegisterProductServiceServer(s, &ss)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())

	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		s.Stop()
	}

	return pb.NewProductServiceClient(conn), closer, nil
}

func server(ctx context.Context) (pb.OrderServiceClient, func(), error) {

	lis := bufconn.Listen(1024 * 1024)

	s := grpc.NewServer()

	h := db.Mock()

	productSvc, _, err := productServer(ctx)

	if err != nil {
		fmt.Println("Failed to listing:", err)
	}

	ss := services.Server{
		H:          h,
		ProductSvc: client.MockInitProductServiceClient(&productSvc),
	}

	pb.RegisterOrderServiceServer(s, &ss)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}), grpc.WithInsecure())

	if err != nil {
		return nil, nil, err
	}

	closer := func() {
		err := lis.Close()
		if err != nil {
			log.Printf("error closing listener: %v", err)
		}
		s.Stop()
	}

	return pb.NewOrderServiceClient(conn), closer, nil
}

func TestOrder(t *testing.T) {
	ctx := context.Background()
	client, closer, err := server(ctx)

	if err != nil {
		t.Fatalf("Failed to setup test server: %v", err)
	}

	defer closer()

	type expectation struct {
		status int64
	}

	tests := map[string]struct {
		req *pb.CreateOrderRequest
		exp expectation
	}{
		"create order": {
			req: &pb.CreateOrderRequest{
				UserId:    1,
				ProductId: 1,
				Quantity:  1,
			},
			exp: expectation{
				status: http.StatusCreated,
			},
		},

		"create order with invalid quantity": {
			req: &pb.CreateOrderRequest{
				UserId:    1,
				ProductId: 2,
				Quantity:  100,
			},
			exp: expectation{
				status: http.StatusConflict,
			},
		},

		"create order with invalid product id": {
			req: &pb.CreateOrderRequest{
				UserId:    1,
				ProductId: 0,
				Quantity:  1,
			},
			exp: expectation{
				status: 404,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := client.CreateOrder(ctx, tc.req)
			if err != nil {
				t.Fatalf("Failed to create order: %v", err)
			}

			if res.Status != tc.exp.status {
				t.Errorf("Expected status %d, got %d", tc.exp.status, res.Status)
			}
		})
	}

}
