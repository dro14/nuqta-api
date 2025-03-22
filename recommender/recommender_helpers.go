package recommender

import "github.com/dro14/nuqta-service/models"

func calculateScore(post *models.Post, now int64) float64 {
	return (2.0*float64(post.Replies) +
		1.5*float64(post.Reposts) +
		1.0*float64(post.Likes) +
		0.5*float64(post.Clicks) +
		0.1*float64(post.Views) -
		1.0*float64(post.Reports)) *
		(2 - float64(now-post.Timestamp)/172800.0)
}
