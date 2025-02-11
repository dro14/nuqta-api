package recommender

import (
	"context"
	"log"
	"slices"
	"time"

	"github.com/dro14/nuqta-service/models"

	"gonum.org/v1/gonum/stat/distuv"
)

func (r *Recommender) UpdateRecs() {
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

		now := time.Now().Unix()
		for i, post := range posts {
			score := calculateScore(post, now)
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
			Alpha: max(rec.Score, 0.00001),
			Beta:  max(float64(rec.Views)-rec.Score, 0.00001),
		}
		recs = append(recs, &models.Post{
			Uid:   rec.Uid,
			Score: beta.Rand(),
		})
	}

	slices.SortFunc(
		recs,
		func(a, b *models.Post) int {
			if a.Score < b.Score {
				return -1
			} else if a.Score > b.Score {
				return 1
			} else {
				return 0
			}
		},
	)

	postUids := make([]string, 0, len(recs))
	for _, rec := range recs {
		postUids = append(postUids, rec.Uid)
	}

	return postUids
}
