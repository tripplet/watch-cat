package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type blockedIP struct {
	InvalidRequests int       `firestore:"invalid_requests"`
	LastInvalid     time.Time `firestore:"last_invalid"`
	BlockedUntil    time.Time `firestore:"blocked_until"`
}

func middlewareIPBlocking(client *firestore.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
		if err != nil {
			return
		}

		entry, err := client.Collection("BlockedIP").Doc(ip).Get(c.Request.Context())
		if err != nil {
			if status.Code(err) == codes.NotFound {
				// IP no blocked, continue processing the request
				return
			}

			log.Println(err)
			return
		}

		var blockedIP blockedIP
		if err := entry.DataTo(&blockedIP); err != nil {
			log.Println(err)
			return
		}

		if (blockedIP.BlockedUntil != time.Time{}) && time.Now().UTC().Before(blockedIP.BlockedUntil) {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
}

func incBlockIP(client *firestore.Client, c *gin.Context) {
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return
	}

	ctx := context.Background()

	dbEntry := client.Collection("BlockedIP").Doc(ip)
	entry, err := dbEntry.Get(ctx)

	entry.Exists()

	if err != nil {
		if status.Code(err) == codes.NotFound {
			_, err := dbEntry.Create(ctx, blockedIP{
				InvalidRequests: 1,
				LastInvalid:     time.Now().UTC(),
				BlockedUntil:    time.Time{}})
			if err != nil {
				log.Println(err)
			}

			return
		}

		log.Println(err)
		return
	}

	var blockedIP blockedIP
	if err := entry.DataTo(&blockedIP); err != nil {
		log.Print(err)
		return
	}

	if blockedIP.InvalidRequests >= 9 {
		entry.Ref.Update(ctx, []firestore.Update{
			{Path: "invalid_requests", Value: blockedIP.InvalidRequests + 1},
			{Path: "blocked_until", Value: time.Now().UTC().Add(1 * time.Hour)},
		})
	} else {
		entry.Ref.Update(ctx, []firestore.Update{{Path: "invalid_requests", Value: blockedIP.InvalidRequests + 1}})
	}
}
