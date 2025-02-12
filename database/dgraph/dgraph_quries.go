package dgraph

const userByUidQuery = `{
	users(func: uid(%s)) {
		uid
		firebase_uid
		name
		username
		bio
		joined_at
		birthday
		banner
		avatars
		thumbnails
		posts: count(~author)
		following: count(follow)
		followers: count(~follow)
	}
}`

const userByFirebaseUidQuery = `{
	users(func: eq(firebase_uid, "%s")) {
		uid
		firebase_uid
		name
		username
		bio
		joined_at
		birthday
		banner
		avatars
		thumbnails
		posts: count(~author)
		following: count(follow)
		followers: count(~follow)
	}
}`

const postQuery = `{
	posts(func: uid(%s)) {
		uid
		text
		posted_at
		author {
			uid
		}
		in_reply_to {
			uid
			text
			posted_at
			author {
				uid
			}
			replies: count(~in_reply_to)
			reposts: count(~repost)
			likes: count(~like)
			views: count(~view)
		}
		replies: count(~in_reply_to)
		reposts: count(~repost)
		likes: count(~like)
		views: count(~view)
	}
}`

const edgeQuery = `{
	edges(func: uid(%s)) {
		%s @filter(uid(%s)) {
			uid
		}
	}
}`

const isRepliedQuery = `{
	edges(func: uid(%s)) {
		~author {
			in_reply_to @filter(uid(%s)) {
				uid
			}
		}
	}
}`

const recentPostsQuery = `{
	posts(func: gt(posted_at, "%d")) @filter(not has(in_reply_to)) {
		uid
		posted_at
		replies: count(~in_reply_to)
		reposts: count(~repost)
		likes: count(~like)
		clicks: count(~click)
		views: count(~view)
		removes: count(~remove)
	}
}`

const followingQuery = `{
	var(func: uid(%s)) {
		following as follow
	}

	posts(func: lt(posted_at, "%d")) @filter(uid_in(author, uid(following)) OR uid_in(~repost, uid(following))) (orderdesc: posted_at, first: 20) {
		uid
	}
}`

const userPostsQuery = `{
	users(func: uid(%s)) {
		posts: ~author @filter(lt(posted_at, "%d") AND not has(in_reply_to)) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`

const userRepliesQuery = `{
	users(func: uid(%s)) {
		posts: ~author @filter(lt(posted_at, "%d") AND has(in_reply_to)) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`

const userRepostsQuery = `{
	users(func: uid(%s)) {
		posts: repost @facets(lt(timestamp, "%d")) @facets(orderdesc: timestamp, first: 20) {
			uid
    	}
	}
}`

const userLikesQuery = `{
	users(func: uid(%s)) {
		posts: like @facets(lt(timestamp, "%d")) @facets(orderdesc: timestamp, first: 20) {
			uid
    	}
	}
}`

const popularRepliesQuery = `{
	posts(func: uid(%s)) {
		replies: ~in_reply_to(orderdesc: val(score), first: 20, offset: %d) {
			uid
			replies: count(~in_reply_to)
			reposts: count(~repost)
			likes: count(~like)
			clicks: count(~click)
			views: count(~view)
			score as math(2.0*replies + 1.5*reposts + likes + 0.5*clicks + 0.1*views)
		}
	}
}`

const recentRepliesQuery = `{
	posts(func: uid(%s)) {
		replies: ~in_reply_to @filter(lt(posted_at, "%d")) (orderdesc: posted_at, first: 20) {
			uid
		}
	}
}`
