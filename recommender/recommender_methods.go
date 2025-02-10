package recommender

import (
	"context"
	"log"
	"slices"
	"time"

	"github.com/dro14/nuqta-service/models"

	"gonum.org/v1/gonum/stat/distuv"
)

func calculateScore(post *models.Post) float64 {
	return 1.0*float64(post.Likes) +
		1.5*float64(post.Reposts) +
		2.0*float64(post.Replies) +
		0.5*float64(post.Clicks) +
		0.1*float64(post.Views) -
		1.0*float64(post.Removes)
}

func (r *Recommender) updateRecs() {
	updateTime, err := r.cache.GetRecUpdateTime()
	if err != nil {
		log.Println(err)
		updateTime = time.Now().Add(time.Minute)
	}
	for {
		time.Sleep(time.Until(updateTime))
		updateTime, err = r.cache.IncrRecUpdateTime()
		if err != nil {
			log.Println(err)
			time.Sleep(time.Minute)
			continue
		}

		posts, err := r.db.GetLatestPosts(context.Background())
		if err != nil {
			log.Println(err)
			time.Sleep(time.Minute)
			continue
		}

		C := 0.0
		m := 0.0

		for i, post := range posts {
			score := calculateScore(post)
			posts[i].Score = score
			C += score
			m += float64(post.Views)
		}

		n := float64(len(posts))
		C /= n
		m /= n

		recs := make([]*models.Post, 0, len(posts))
		for _, post := range posts {
			v := float64(post.Views)
			recs = append(recs, &models.Post{
				Uid:   post.Uid,
				Score: (v*post.Score + m*C) / (v + m),
				Views: post.Views,
			})
		}

		r.recs = recs
		err = r.cache.SetRecs(recs)
		if err != nil {
			log.Println(err)
		}
	}
}

func (r *Recommender) GetRecs() []string {
	recs := make([]*models.Post, 0, len(r.recs))
	for _, rec := range r.recs {
		beta := distuv.Beta{
			Alpha: rec.Score,
			Beta:  float64(rec.Views) - rec.Score,
		}
		recs = append(recs, &models.Post{
			Uid:   rec.Uid,
			Score: beta.Rand(),
		})
	}
	slices.SortFunc(recs, func(a, b *models.Post) int {
		if a.Score < b.Score {
			return -1
		} else if a.Score > b.Score {
			return 1
		} else {
			return 0
		}
	})

	postUids := make([]string, 0, len(recs))
	for _, rec := range recs {
		postUids = append(postUids, rec.Uid)
	}

	return postUids
}
