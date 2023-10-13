package lib

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/joho/godotenv"
	"golang.ngrok.com/ngrok"
	"golang.ngrok.com/ngrok/config"
	"golang.org/x/sync/errgroup"
)

func RunTunnel(port *int) error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("error loading env variables: %s", err)
	}
	ctx := context.Background()
	l, err := ngrok.Listen(ctx,
		config.TCPEndpoint(),
		ngrok.WithAuthtoken(os.Getenv("NGROK_AUTHTOKEN")),
	)
	if err != nil {
		return err
	}
	fmt.Println("IP: ", l.URL())
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go func() {
			err := handleConn(ctx, port, conn)
			fmt.Println("Connection closed: ", err)
		}()
	}
}

func handleConn(ctx context.Context, port *int, conn net.Conn) error {
	next, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return err
	}

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		_, err := io.Copy(next, conn)
		return err
	})
	g.Go(func() error {
		_, err := io.Copy(conn, next)
		return err
	})

	return g.Wait()

}
