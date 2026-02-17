package services

import (
	"blog/internal/cache"
	repo "blog/internal/repositories"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
	"golang.org/x/sync/singleflight"
)

// StrapiService contains all methods for interacting with Strapi
type StrapiService struct {
	client graphql.Client
	token  string              // Authorization bearer token
	sf     *singleflight.Group // Singleflight group to reduce api calls
	cache  *cache.Cache        // Redis cache
}

func NewStrapiService(endpoint, token string, cache *cache.Cache) *StrapiService {
	httpClient := &http.Client{
		Transport: &authTransport{
			token: token,
			base:  http.DefaultTransport,
		},
	}

	return &StrapiService{
		client: graphql.NewClient(endpoint, httpClient),
		token:  token,
		sf:     &singleflight.Group{},
		cache:  cache,
	}
}

// GetAllPosts returns all posts from Strapi
func (s *StrapiService) GetAllPosts(ctx context.Context, pageSize, page int) ([]repo.Post, error) {
	key := fmt.Sprintf("strapi:posts:all:%d:%d", pageSize, page)

	// Check if exists in cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var posts []repo.Post
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			return posts, nil
		}
	}

	//  Cache miss - use singleflight
	result, err, _ := s.sf.Do(key, func() (any, error) {
		resp, err := repo.GetAllPosts(ctx, s.client, pageSize, page)
		if err != nil {
			return nil, err
		}

		// Store in cache
		go func() {
			if data, err := json.Marshal(resp.Posts); err == nil {
				// Set cache with 5 minute expiry
				s.cache.Set(context.Background(), key, data, time.Minute*5)
			}
		}()

		return resp.Posts, nil
	})

	if err != nil {
		return nil, err
	}

	// Return result
	posts, ok := result.([]repo.Post)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight")
	}
	return posts, nil
}

// GetFeaturedPosts returns any post that has "featured" as true from Strapi
func (s *StrapiService) GetFeaturedPosts(ctx context.Context, pageSize, page int) ([]repo.PostSummary, error) {
	key := fmt.Sprintf("strapi:posts:featured:%d:%d", pageSize, page)

	// Check if exists in cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var posts []repo.PostSummary
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			return posts, nil
		}
	}

	// Cache miss - use singleflight
	result, err, _ := s.sf.Do(key, func() (any, error) {
		resp, err := repo.GetFeaturedPosts(ctx, s.client, pageSize, page)
		if err != nil {
			return nil, err
		}

		// Store in cache
		go func() {
			if data, err := json.Marshal(resp.Posts); err == nil {
				// Set cache with 5 minute expiry
				s.cache.Set(context.Background(), key, data, time.Minute*5)
			}
		}()

		return resp.Posts, nil
	})

	if err != nil {
		return nil, err
	}

	// Return result
	posts, ok := result.([]repo.PostSummary)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight")
	}
	return posts, nil
}

// GetPost returns a specific post from Strapi
func (s *StrapiService) GetPost(ctx context.Context, slug string) (*repo.Post, error) {
	key := fmt.Sprintf("strapi:post:%s", slug)

	// Check if exists in cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var post repo.Post
		if err := json.Unmarshal([]byte(cached), &post); err == nil {
			return &post, nil
		}
	}

	// Cache miss - use singleflight
	result, err, _ := s.sf.Do(key, func() (any, error) {
		resp, err := repo.GetPost(ctx, s.client, slug)
		if err != nil {
			return nil, err
		}

		if len(resp.Posts) == 0 {
			return nil, fmt.Errorf("post not found: %s", slug)
		}

		post := &resp.Posts[0] // Create pointer

		// Store in cache
		go func() {
			if data, err := json.Marshal(post); err == nil {
				// Set cache with 15 minute expiry
				s.cache.Set(context.Background(), key, data, time.Minute*15)
			}
		}()

		return post, nil
	})

	if err != nil {
		return nil, err
	}

	// Return result
	post, ok := result.(*repo.Post)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight")
	}
	return post, nil
}

// GetPostSummaries returns post summaries from Strapi
func (s *StrapiService) GetPostSummaries(ctx context.Context, pageSize, page int) ([]repo.PostSummary, error) {
	key := fmt.Sprintf("strapi:posts:summary:%d:%d", pageSize, page)

	// Check if exists in cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var posts []repo.PostSummary
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			return posts, nil
		}
	}

	// Cache miss - use singleflight
	result, err, _ := s.sf.Do(key, func() (any, error) {
		resp, err := repo.GetPostSummaries(ctx, s.client, pageSize, page)
		if err != nil {
			return nil, err
		}

		// Store in cache
		go func() {
			if data, err := json.Marshal(resp.Posts); err == nil {
				// Set cache with 5 minute expiry
				s.cache.Set(context.Background(), key, data, time.Minute*5)
			}
		}()

		return resp.Posts, nil
	})

	if err != nil {
		return nil, err
	}

	// Return results
	posts, ok := result.([]repo.PostSummary)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight")
	}
	return posts, nil
}

// GetPostsByTag returns posts with a specific tag from Strapi
func (s *StrapiService) GetPostsByTag(ctx context.Context, tag string, pageSize, page int) ([]repo.PostSummary, error) {
	key := fmt.Sprintf("strapi:posts:tag:%s:%d:%d", tag, pageSize, page)

	// Check if exists in cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var posts []repo.PostSummary
		if err := json.Unmarshal([]byte(cached), &posts); err == nil {
			return posts, nil
		}
	}

	// Cache miss - use singleflight
	result, err, _ := s.sf.Do(key, func() (any, error) {
		resp, err := repo.GetPostsByTag(ctx, s.client, tag, pageSize, page)
		if err != nil {
			return nil, err
		}

		// Store in cache
		go func() {
			if data, err := json.Marshal(resp.Posts); err == nil {
				// Store in cache with 5 minute expiry
				s.cache.Set(context.Background(), key, data, time.Minute*5)
			}
		}()

		return resp.Posts, nil
	})

	if err != nil {
		return nil, err
	}

	// Return results
	posts, ok := result.([]repo.PostSummary)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight")
	}
	return posts, nil
}

// GetAllAuthors returns all authors from Strapi
func (s *StrapiService) GetAllAuthors(ctx context.Context, pageSize, page int) ([]repo.Author, error) {
	key := fmt.Sprintf("strapi:authors:all:%d:%d", pageSize, page)

	// Check if exists in cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var authors []repo.Author
		if err := json.Unmarshal([]byte(cached), &authors); err == nil {
			return authors, nil
		}
	}

	// Cache miss - use singleflight
	result, err, _ := s.sf.Do(key, func() (any, error) {
		resp, err := repo.GetAllAuthors(ctx, s.client, pageSize, page)
		if err != nil {
			return nil, err
		}

		// Store in cache
		go func() {
			if data, err := json.Marshal(resp.Authors); err == nil {
				// Store in cache with 5 minute expiry
				s.cache.Set(context.Background(), key, data, time.Minute*5)
			}
		}()

		return resp.Authors, nil
	})

	if err != nil {
		return nil, err
	}

	// Return results
	authors, ok := result.([]repo.Author)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight")
	}
	return authors, nil
}

// GetAllTags returns all tags from Strapi
func (s *StrapiService) GetAllTags(ctx context.Context) ([]repo.Tag, error) {
	key := fmt.Sprintf("strapi:tags:all")

	// Check if exists in cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil {
		var tags []repo.Tag
		if err := json.Unmarshal([]byte(cached), &tags); err == nil {
			return tags, nil
		}
	}

	// Cache miss - use singleflight
	result, err, _ := s.sf.Do(key, func() (any, error) {
		resp, err := repo.GetAllTags(ctx, s.client)
		if err != nil {
			return nil, err
		}

		// Store in cache
		go func() {
			if data, err := json.Marshal(resp.Tags); err == nil {
				// Store in cache with 5 minute expiry
				s.cache.Set(context.Background(), key, data, time.Minute*5)
			}
		}()

		return resp.Tags, nil
	})

	if err != nil {
		return nil, err
	}

	// Return results
	tags, ok := result.([]repo.Tag)
	if !ok {
		return nil, fmt.Errorf("unexpected type from singleflight")
	}
	return tags, nil
}

// Auth transport for headers
type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Bearer "+t.token)
	return t.base.RoundTrip(req)
}
