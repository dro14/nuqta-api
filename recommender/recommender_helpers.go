package recommender

import "github.com/dro14/nuqta-service/models"

func calculateScore(post *models.Post, now int64) float64 {
	return (1.0*float64(post.Likes) +
		1.5*float64(post.Reposts) +
		2.0*float64(post.Replies) +
		0.5*float64(post.Clicks) +
		0.1*float64(post.Views) -
		1.0*float64(post.Removes)) *
		(2 - float64(now-post.PostedAt)/172800.0)
}
